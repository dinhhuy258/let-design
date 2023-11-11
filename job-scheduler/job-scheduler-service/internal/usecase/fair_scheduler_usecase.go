package usecase

import (
	"container/list"
	"container/ring"
	"context"
	"job-scheduler-service/internal/entity"
	"sort"

	"github.com/samber/lo"
)

type fairSchedulerUsecase struct {
	baseSchedulerUsecase
}

func NewFairSchedulerUsecase() SchedulerUsecase {
	return &fairSchedulerUsecase{}
}

func (_self *fairSchedulerUsecase) ScheduleJobs(ctx context.Context, jobs []entity.Job) error {
	jobGroups := lo.PartitionBy(jobs, func(job entity.Job) int {
		return int(job.User.Id)
	})

	// sort job groups by user's job weight
	sort.Slice(jobGroups, func(i, j int) bool {
		return jobGroups[i][0].User.JobWeight > jobGroups[j][0].User.JobWeight
	})

	for jobGroupIdx := range jobGroups {
		// sort jobs by execute_at
		sort.Slice(jobGroups[jobGroupIdx], func(i, j int) bool {
			return jobGroups[jobGroupIdx][i].ExecuteAt.Before(jobGroups[jobGroupIdx][j].ExecuteAt)
		})
	}

	scheduledJobs := _self.getScheduledJobs(ctx, jobs)
	for _, job := range scheduledJobs {
		_self.scheduleJob(ctx, job)
	}

	return nil
}

func (*fairSchedulerUsecase) getScheduledJobs(ctx context.Context, jobs []entity.Job) []entity.Job {
	// partition jobs by user id
	jobGroups := lo.PartitionBy(jobs, func(job entity.Job) int {
		return int(job.User.Id)
	})

	// sort job groups by user's job weight
	sort.Slice(jobGroups, func(i, j int) bool {
		return jobGroups[i][0].User.JobWeight > jobGroups[j][0].User.JobWeight
	})

	for jobGroupIdx := range jobGroups {
		// sort jobs by execute_at
		sort.Slice(jobGroups[jobGroupIdx], func(i, j int) bool {
			return jobGroups[jobGroupIdx][i].ExecuteAt.Before(jobGroups[jobGroupIdx][j].ExecuteAt)
		})
	}

	// userWeightsMap keeps track of the user's weight
	userWeightsMap := make(map[uint64]float32)
	// create a circular linked list to keep track of the job groups
	jobGroupsRing := ring.New(len(jobGroups))
	// initialize the ring and set the value of each element to the list of jobs of each job group
	for i := 0; i < len(jobGroups); i++ {
		// keep track of jobGroups in a queue
		// it will help us to remove the processed tasks from the queue easilys
		// we can still use slice, but we need to sort the slice in descending order
		// and remove the processed tasks from the slice from the end of the slice to the beginning of
		// the slice (it will be O(1) time complexity for deletion operation)
		// but using queue is more readable
		jobGroupsRing.Value = list.New()
		for _, job := range jobGroups[i] {
			jobGroupsRing.Value.(*list.List).PushBack(job)
		}

		// represents the user of the job group
		user := jobGroups[i][0].User
		userWeightsMap[user.Id] = user.JobWeight

		jobGroupsRing = jobGroupsRing.Next()
	}

	scheduledJobs := make([]entity.Job, 0, len(jobs))

	// fair scheduling
	for jobGroupsRing.Len() > 1 {
		currentGroup := jobGroupsRing.Value.(*list.List)
		user := currentGroup.Front().Value.(entity.Job).User

		for userWeightsMap[user.Id] > 1.0 {
			userWeightsMap[user.Id] = userWeightsMap[user.Id] - 1.0
			// schedule the current job
			scheduledJobs = append(scheduledJobs, currentGroup.Front().Value.(entity.Job))
			// remove the current job from the current job group
			currentGroup.Remove(currentGroup.Front())
			// if the current job group is empty, remove the current job group from the queue
			if currentGroup.Len() == 0 {
				jobGroupsRing = jobGroupsRing.Prev()
				jobGroupsRing.Unlink(1)
				jobGroupsRing = jobGroupsRing.Next()
			}
		}

		userWeightsMap[user.Id] = userWeightsMap[user.Id] + user.JobWeight

		// move to the next job group
		jobGroupsRing = jobGroupsRing.Next()
	}

	// schedule the remaining jobs
	remainingJobGroup := jobGroupsRing.Value.(*list.List)
	for remainingJobGroup.Len() > 0 {
		scheduledJobs = append(scheduledJobs, remainingJobGroup.Front().Value.(entity.Job))
		remainingJobGroup.Remove(remainingJobGroup.Front())
	}

	return scheduledJobs
}
