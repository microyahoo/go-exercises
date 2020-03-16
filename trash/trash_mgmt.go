package trash

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/juju/errors"

	"xsky-demon/log"
	"xsky-demon/models"
)

var (
	// TrashActionRecycle stands for "recycle" operation
	TrashActionRecycle = "recycle"
	// TrashActionRestore stands for "restore" operation
	TrashActionRestore = "restore"
	// TrashActionDelete stands for "delete" operation
	TrashActionDelete = "delete"
	// TrashActionUpdatePolicy stands for "update policy" operation
	TrashActionUpdatePolicy = "update_policy"
)

// Scheduler is a interface for job management.
type Scheduler interface {
	AddJob(job Job)
	RemoveJob(job Job)
	// StopJob(job Job)

	// stop the scheduler
	Stop()
	// run the scheduler
	Run()

	EventChan() chan *Event
	HandleEvent(event *Event) (err error)
}

// Job is a interface for task
type Job interface {
	GetUniqueKey() string
	GetExpireTime() time.Time
	SetScheduler(Scheduler) error
}

// trashMgmt is trash resource management, which implements Scheduler interface.
type trashMgmt struct {
	ttlKeyHeap  *ttlKeyHeap
	waitJobsNum int64

	// the key is combination of resource_type and resource_id
	jobs map[string]Job

	m *models.DemonModel

	lock sync.RWMutex

	stopChan  chan struct{}
	eventChan chan *Event
}

var _ Scheduler = (*trashMgmt)(nil)
var _ Job = (*job)(nil)
var _ Job = (*trashJob)(nil)

// job is a timed task
type job struct {
	// ID         int64
	UniqueKey  string
	ExpireTime time.Time

	// ResourceID   int64
	// ResourceType string

	running bool

	// stopChan chan struct{}
	// eventChan chan *Event

	jobFunc string
	funcs   map[string]interface{}
	fparams map[string][]interface{}

	mgmt Scheduler
}

func (j *job) GetUniqueKey() string {
	// return strconv.FormatInt(j.ID, 10)
	return j.UniqueKey
}

func (j *job) SetScheduler(s Scheduler) error {
	j.mgmt = s
	return nil
}

func (j *job) GetExpireTime() time.Time {
	return j.ExpireTime
}

type trashJob struct {
	job

	ResourceID   int64
	ResourceType string
}

// Event defines the action on which job will take.
type Event struct {
	Action string
	Status string

	Resource *models.TrashResource
	Policy   *models.TrashPolicy
}

func (e *Event) String() string {
	return fmt.Sprintf("Event[Action=%s, Status=%s, Resource[Type=%s, ID=%d]]",
		e.Action, e.Status,
		e.Resource.ResourceType,
		e.Resource.ResourceID)
}

func (j *job) init(expireTime time.Time) {
	j.ExpireTime = expireTime
	j.running = false
	// j.stopChan = make(chan struct{}, 0)
	j.funcs = make(map[string]interface{})
	j.fparams = make(map[string][]interface{})
}

// AddFunc registers callback
func (j *job) AddFunc(jobFun interface{}, params ...interface{}) {
	// TODO
}

// newTrashJob returns a new trash job.
func (j *trashJob) newTrashJob(resourceType string, resourceID int64, expireTime time.Time) {
	j.init(expireTime)
	j.ResourceID = resourceID
	j.ResourceType = resourceType
	j.UniqueKey = keysHelper(j.ResourceType, strconv.FormatInt(j.ResourceID, 10))
}

var (
	singletonMgmt Scheduler
	once          sync.Once
)

// NewScheduler returns a new trash management.
func NewScheduler() (s Scheduler, err error) {
	once.Do(func() {
		var m *models.DemonModel
		m, err = models.NewModel()
		if err != nil {
			return
		}
		singletonMgmt = &trashMgmt{
			ttlKeyHeap:  newTTLKeyHeap(),
			jobs:        make(map[string]Job),
			m:           m,
			waitJobsNum: 0,
			stopChan:    make(chan struct{}),
			eventChan:   make(chan *Event, 0),
			// workerMap:   make(map[string]*TrashWorker),
		}
	})
	if err != nil {
		return s, errors.Trace(err)
	}
	return singletonMgmt, nil
}

func (mgmt *trashMgmt) EventChan() chan *Event {
	return mgmt.eventChan
}

func (mgmt *trashMgmt) Run() {
	var (
		duration time.Duration
		timer    *time.Timer
		j        Job
		dumyJob  = new(job)
		now      = time.Now()
		err      error
	)
	dumyJob.init(now.Add(math.MaxInt64))
	mgmt.AddJob(dumyJob)
	log.Infof("Start to run trash scheduler")
	for {
		if mgmt.hasTTLKeys() {
			j = mgmt.ttlKeyHeap.top()
			duration = j.GetExpireTime().Sub(time.Now())
			// if duration is negative, the timer will fire immediately
			timer = time.NewTimer(duration)
		} else {
			log.Errorf("The ttl key heap should not be null.")
			return
		}
		for {
			select {
			case now = <-timer.C:
				mgmt.RemoveJob(j)
				// timer.Stop()
				// TODO (zhengliang): run the task to delete volume
				if j, ok := j.(*trashJob); ok {
					e := &Event{
						Action: TrashActionDelete,
						Status: models.StatusActive,
						Resource: &models.TrashResource{
							ResourceID:   j.ResourceID,
							ResourceType: j.ResourceType,
							ExpireTime:   now,
							Status:       models.StatusActive,
							ActionStatus: models.StatusActive,
						},
					}
					log.Infof("Before delete job, send a event %v", e)
					go func() { mgmt.eventChan <- e }()
					log.Info("After send a event %v", e)
				}
			case event := <-mgmt.eventChan:
				log.Infof("Received the event: %v", event)
				timer.Stop()
				err = mgmt.HandleEvent(event)
				if err != nil {
					log.Errorf("Failed to handle event %v: %v", event, err)
					// TODO (zhengliang): handle the error
				}
			case <-mgmt.stopChan:
				log.Infof("Stop trash management.")
				timer.Stop()
				return
			}
			break
		}
	}
}

// HandleEvent dispatches or handles events
func (mgmt *trashMgmt) HandleEvent(event *Event) (err error) {
	if event == nil {
		return nil
	}
	log.Infof("workerMap: %v", workerMap)
	if event.Resource != nil {
		worker, ok := workerMap[event.Resource.ResourceType]
		if !ok {
			return errors.Errorf("Doesn't support %s", event.Resource.ResourceType)
		}
		switch event.Action {
		case TrashActionRecycle:
			log.Infof("Create a new trash job %v, trashMgmt: %v", event, mgmt)
			j := new(trashJob)
			j.newTrashJob(event.Resource.ResourceType,
				event.Resource.ResourceID, event.Resource.ExpireTime)
			mgmt.AddJob(j)
			return nil
			// return worker.Recycle(event.Resource)
		case TrashActionRestore:
			return worker.Restore(event.Resource)
		case TrashActionDelete:
			log.Infof("Receive delete action: %v, trashMgmt: %v", event, mgmt)
			return worker.Delete(event.Resource)
		default:
			return errors.Errorf("Unknown action type %s", event.Action)
		}
	}
	if event.Policy != nil {
		worker, ok := workerMap[event.Policy.ResourceType]
		if !ok {
			return errors.Errorf("Doesn't support %s", event.Policy.ResourceType)
		}
		switch event.Action {
		case TrashActionRecycle:
			return worker.UpdatePolicy(event.Policy)
		default:
			return errors.Errorf("Unknown action type %s", event.Action)
		}
	}
	return nil
}

func (mgmt *trashMgmt) loadDataFromDB() {
}

// StopJob stop specified job.
// func (mgmt *trashMgmt) StopJob(j Job) {
// 	if j == nil {
// 		log.Warnf("The job is null")
// 		return
// 	}
// 	mgmt.RemoveJob(j)
// 	// TODO (zhengliang): need to confirm
// 	// j.stopChan <- struct{}{}
// }

// Stop stops trash management.
func (mgmt *trashMgmt) Stop() {
	mgmt.stopChan <- struct{}{}
}

// AddJob add job to scheduler.
func (mgmt *trashMgmt) AddJob(j Job) {
	if j == nil {
		log.Warnf("The job is null")
		return
	}
	mgmt.lock.Lock()
	defer mgmt.lock.Unlock()
	mgmt.ttlKeyHeap.push(j)
	mgmt.waitJobsNum++
	mgmt.jobs[j.GetUniqueKey()] = j
	j.SetScheduler(mgmt)
}

// RemoveJob removes job from scheduler
func (mgmt *trashMgmt) RemoveJob(j Job) {
	if j == nil {
		log.Warnf("The job is null")
		return
	}
	mgmt.lock.Lock()
	defer mgmt.lock.Unlock()
	uniqueKey := j.GetUniqueKey()
	if j, ok := mgmt.jobs[uniqueKey]; ok {
		mgmt.ttlKeyHeap.remove(j)
		mgmt.waitJobsNum--
		delete(mgmt.jobs, uniqueKey)
	} else {
		log.Warnf("The job was not found.")
	}
}

func (mgmt *trashMgmt) getJob(uniqueKey string) Job {
	if _, ok := mgmt.jobs[uniqueKey]; !ok {
		return nil
	}
	return mgmt.jobs[uniqueKey]
}

func (mgmt *trashMgmt) getJobOfTrashResource(res *models.TrashResource) Job {
	return mgmt.getJob(keysHelper(res.ResourceType, strconv.FormatInt(res.ResourceID, 10)))
}

func (mgmt *trashMgmt) hasTTLKeys() bool {
	mgmt.lock.Lock()
	defer mgmt.lock.Unlock()
	return mgmt.ttlKeyHeap.Len() != 0
}

func keysHelper(keys ...string) string {
	return strings.Join(keys, "-")
}

// Worker is a interface which includes some trash operations
type Worker interface {
	Recycle(res *models.TrashResource) error
	Restore(res *models.TrashResource) error
	Delete(res *models.TrashResource) error
	UpdatePolicy(policy *models.TrashPolicy) error
}

// the key is resource type
var (
	workerMap map[string]Worker
)

// RegisterWorker register worker
func RegisterWorker(resourceType string, worker Worker) {
	if workerMap == nil {
		workerMap = make(map[string]Worker)
	}
	workerMap[resourceType] = worker
	// TODO need to confirm
	// if singletonMgmt == nil {
	// 	newTrashMgmt()
	// }
	// worker.mgmt = singletonMgmt
}
