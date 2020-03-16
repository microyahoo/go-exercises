package main

import (
	"fmt"
	"time"
)

type Job interface {
}

type job struct {
	ExpireTime time.Time
	UniqueKey  string
}

func (j *job) init(expireTime time.Time) {
	j.ExpireTime = expireTime
}

type trashJob struct {
	job
	ResourceType string
	ResourceID   int64
}

func (j *trashJob) newTrashJob(resourceType string, resourceID int64, expireTime time.Time) {
	j.init(expireTime)
	j.ResourceType = resourceType
	j.ResourceID = resourceID
}

func (j *trashJob) String() string {
	return fmt.Sprintf("ResourceType = %s, ResourceID = %d, ExpireTime = %v", j.ResourceType, j.ResourceID, j.ExpireTime)
}

func doJob(job Job) {
	// fmt.Println(job.expireTime)
	// fmt.Println(job.uniqueKey)
	// switch t := interface{}(job).(type) {
	// case trashJob:
	// 	fmt.Println(t)
	// }
	if job, ok := job.(*trashJob); ok {
		fmt.Println(job)
	} else {
		fmt.Println("hello")
	}
}

func main() {
	j := new(trashJob)
	j.newTrashJob("block", 1, time.Now())
	doJob(j)

	ticker := time.NewTicker(1 * time.Second)
	ticker.Stop()
	ticker.Stop()
}
