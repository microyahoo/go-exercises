package main

import (
	"context"
	"encoding/json"
	"time"

	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers/internalinterfaces"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

const (
	Group                           string = "argoproj.io"
	Version                         string = "v1alpha1"
	APIVersion                      string = Group + "/" + Version
	WorkflowKind                    string = "Workflow"
	WorkflowSingular                string = "workflow"
	WorkflowPlural                  string = "workflows"
	WorkflowFullName                string = WorkflowPlural + "." + Group
	LabelKeyControllerInstanceID           = WorkflowFullName + "/controller-instanceid"
	LabelKeyWorkflowArchivingStatus        = WorkflowFullName + "/workflow-archiving-status"
)

var (
	instanceID = ""
	ns         = "argo-data-closeloop"
	kubeconfig = "/tmp/config"
	// kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
)

func muste(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	muste(err)

	dynamicInterface, err := dynamic.NewForConfig(restConfig)
	muste(err)
	argoInformer := newWorkflowInformer(dynamicInterface, ns, 20*time.Minute, tweakListRequestListOptions, tweakWatchRequestListOptions, nil)

	// 添加Event事件处理函数
	_, err = argoInformer.AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: func(obj interface{}) bool {
			un, ok := obj.(*unstructured.Unstructured)
			// no need to check the `common.LabelKeyCompleted` as we already know it must be complete
			return ok && un.GetLabels()[LabelKeyWorkflowArchivingStatus] == "Pending"
		},
		Handler: cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				muste(err)
				klog.Infof("receive add event, key=%s", key)
				un, ok := obj.(*unstructured.Unstructured)
				if !ok {
					return
				}
				wf, err := fromUnstructured(un)
				muste(err)
				klog.Infof("receive add event: %s/%s", wf.Namespace, wf.Name)
			},
			UpdateFunc: func(_, obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				muste(err)
				klog.Infof("receive update event, key=%s", key)
				un, ok := obj.(*unstructured.Unstructured)
				if !ok {
					return
				}
				wf, err := fromUnstructured(un)
				muste(err)
				klog.Infof("receive update event: %s/%s", wf.Namespace, wf.Name)
			},
		},
	})
	muste(err)

	stopCh := make(chan struct{})
	defer close(stopCh)
	// 启动 informer
	go argoInformer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, argoInformer.HasSynced) {
		panic("Timed out waiting for caches to sync")
	}

	<-stopCh
}

func newWorkflowInformer(dclient dynamic.Interface, ns string, resyncPeriod time.Duration, tweakListRequestListOptions internalinterfaces.TweakListOptionsFunc, tweakWatchRequestListOptions internalinterfaces.TweakListOptionsFunc, indexers cache.Indexers) cache.SharedIndexInformer {
	resource := schema.GroupVersionResource{
		Group:    Group,
		Version:  "v1alpha1",
		Resource: WorkflowPlural,
	}
	informer := newFilteredUnstructuredInformer(
		resource,
		dclient,
		ns,
		resyncPeriod,
		indexers,
		tweakListRequestListOptions,
		tweakWatchRequestListOptions,
	)
	return informer
}

func newFilteredUnstructuredInformer(resource schema.GroupVersionResource, client dynamic.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListRequestListOptions internalinterfaces.TweakListOptionsFunc, tweakWatchRequestListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	ctx := context.Background()
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListRequestListOptions != nil {
					tweakListRequestListOptions(&options)
				}
				return client.Resource(resource).Namespace(namespace).List(ctx, options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakWatchRequestListOptions != nil {
					tweakWatchRequestListOptions(&options)
				}
				return client.Resource(resource).Namespace(namespace).Watch(ctx, options)
			},
		},
		&unstructured.Unstructured{},
		resyncPeriod,
		indexers,
	)
}

func instanceIDRequirement(instanceID string) labels.Requirement {
	var instanceIDReq *labels.Requirement
	var err error
	if instanceID != "" {
		instanceIDReq, err = labels.NewRequirement(LabelKeyControllerInstanceID, selection.Equals, []string{instanceID})
	} else {
		instanceIDReq, err = labels.NewRequirement(LabelKeyControllerInstanceID, selection.DoesNotExist, nil)
	}
	if err != nil {
		panic(err)
	}
	return *instanceIDReq
}

func tweakListRequestListOptions(options *metav1.ListOptions) {
	labelSelector := labels.NewSelector().
		Add(instanceIDRequirement(instanceID))
	options.LabelSelector = labelSelector.String()
	// `ResourceVersion=0` does not honor the `limit` in API calls, which results in making significant List calls
	// without `limit`. For details, see https://github.com/argoproj/argo-workflows/pull/11343
	options.ResourceVersion = ""
}

func tweakWatchRequestListOptions(options *metav1.ListOptions) {
	labelSelector := labels.NewSelector().
		Add(instanceIDRequirement(instanceID))
	options.LabelSelector = labelSelector.String()
}

func fromUnstructured(un *unstructured.Unstructured) (*wfv1.Workflow, error) {
	var wf wfv1.Workflow
	err := fromUnstructuredObj(un, &wf)
	return &wf, err
}

func fromUnstructuredObj(un *unstructured.Unstructured, v interface{}) error {
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(un.Object, v)
	if err != nil {
		if err.Error() == "cannot convert int64 to v1alpha1.AnyString" {
			data, err := json.Marshal(un)
			if err != nil {
				return err
			}
			return json.Unmarshal(data, v)
		}
		return err
	}
	return nil
}
