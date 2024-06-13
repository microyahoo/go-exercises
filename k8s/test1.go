package main

import (
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func Muste(e error) {
	if e != nil {
		panic(e)
	}
}
func InitClientSet() *kubernetes.Clientset {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	Muste(err)
	return kubernetes.NewForConfigOrDie(restConfig)
}

func main() {
	clientSet := InitClientSet()
	// 初始化 sharedInformerFactory
	sharedInformerFactory := informers.NewSharedInformerFactory(clientSet, 0)
	// 生成podInformer
	podInformer := sharedInformerFactory.Core().V1().Pods()
	// 生成具体 informer/indexer
	informer := podInformer.Informer()
	indexer := podInformer.Lister()
	// 添加Event事件处理函数
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			podObj := obj.(*corev1.Pod)

			fmt.Printf("AddFunc: %s\n", podObj.GetName())

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldPodObj := oldObj.(*corev1.Pod)
			newPodObj := newObj.(*corev1.Pod)

			fmt.Printf("old: %s\n", oldPodObj.GetName())
			fmt.Printf("new: %s\n", newPodObj.GetName())

		},
		DeleteFunc: func(obj interface{}) {
			podObj := obj.(*corev1.Pod)
			fmt.Printf("deleteFunc: %s\n", podObj.GetName())
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)
	// 启动 informer
	sharedInformerFactory.Start(stopCh)
	// 等待同步完成
	sharedInformerFactory.WaitForCacheSync(stopCh)

	// 利用 indexer 获取资源
	pods, err := indexer.List(labels.Everything())
	Muste(err)
	for _, items := range pods {
		fmt.Printf("namespace: %s, name:%s\n", items.GetNamespace(), items.GetName())
	}
	<-stopCh
}

// 链接：https://juejin.cn/post/7082728781458178085
