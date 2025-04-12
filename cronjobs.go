package cronjobsmanager

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type CronJobsManager struct {
	jobs map[string]*JobInfo
	mu   sync.RWMutex
}

type JobInfo struct {
	Id          string
	Description string
	context     context.Context
	Cancel      context.CancelFunc
}

type Schedule struct {
	Every   time.Duration
	StartAt time.Time
}

func NewCronJobsManager() *CronJobsManager {
	return &CronJobsManager{
		jobs: make(map[string]*JobInfo),
		mu:   sync.RWMutex{},
	}
}

func (x *CronJobsManager) createJob(id string, description string, ctx context.Context) *JobInfo {
	ctx, cancelFunc := context.WithCancel(ctx)
	return &JobInfo{
		Id:          id,
		Description: description,
		context:     ctx,
		Cancel:      cancelFunc,
	}
}

func (x *CronJobsManager) RemoveJob(jobId string) error {
	if job, ok := x.jobs[jobId]; ok {
		job.Cancel()

		x.mu.Lock()
		delete(x.jobs, jobId)
		x.mu.Unlock()

		return nil
	}

	return errors.New("job not found")
}

func (x *CronJobsManager) ListJobs() []*JobInfo {
	jobList := make([]*JobInfo, 0, len(x.jobs))

	x.mu.RLock()
	for _, value := range x.jobs {
		jobList = append(jobList, value)
	}
	x.mu.RUnlock()

	return jobList
}

func (x *CronJobsManager) RunJob(ctx context.Context, description string, schedule Schedule, jobFunc func()) (jobId string, err error) {
	if schedule.StartAt.IsZero() && schedule.Every == 0 {
		return "", errors.New("invalid schedule")
	}

	id, _ := uuid.NewUUID()
	job := x.createJob(id.String(), description, ctx)

	x.mu.Lock()
	x.jobs[id.String()] = job
	x.mu.Unlock()

	if !schedule.StartAt.IsZero() || schedule.Every > 0{
		go x.handleJob(job, &schedule, jobFunc)
	}

	return job.Id, nil
}

func (x *CronJobsManager) handleJob(job *JobInfo, schedule *Schedule, jobFunc func()){
	if !schedule.StartAt.IsZero(){
		select {
		case <-job.context.Done():
			return
		case <-time.After(time.Until(schedule.StartAt)):
			jobFunc()
		}
	}

	if schedule.Every == 0{
		x.RemoveJob(job.Id)
	}else{
		defer x.RemoveJob(job.Id)

		ticker := time.NewTicker(schedule.Every)
		defer ticker.Stop()

		for {
			select {
			case <-job.context.Done():
				return
			case <-ticker.C:
				jobFunc()
			}
		}
	}
}
