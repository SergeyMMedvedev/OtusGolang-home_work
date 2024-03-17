package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func fillTasksList(tasksCount int, sumTime *time.Duration, tasks *[]Task, runTasksCount *int32) {
	for i := 0; i < tasksCount; i++ {
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
		*sumTime += taskSleep

		*tasks = append(*tasks, func() error {
			time.Sleep(taskSleep)
			atomic.AddInt32(runTasksCount, 1)
			return nil
		})
	}
}

func fillTasksListWithErrors(tasksCount int, tasks *[]Task, runTasksCount *int32) {
	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		*tasks = append(*tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(runTasksCount, 1)
			return err
		})
	}
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		fillTasksListWithErrors(tasksCount, &tasks, &runTasksCount)

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("if were errors in tasks, and M is zero, than finished not more N tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		fillTasksListWithErrors(tasksCount, &tasks, &runTasksCount)
		workersCount := 10
		maxErrorsCount := 0
		err := Run(tasks, workersCount, maxErrorsCount)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("if were no errors in tasks, and M is zero, than finished all tasks", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		var sumTime time.Duration

		fillTasksList(tasksCount, &sumTime, &tasks, &runTasksCount)
		workersCount := 5
		maxErrorsCount := 0
		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		fillTasksList(tasksCount, &sumTime, &tasks, &runTasksCount)

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("the number of tasks is less than the number of workers", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		fillTasksList(tasksCount, &sumTime, &tasks, &runTasksCount)
		workersCount := 20
		maxErrorsCount := 5

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("mixed tasks, but errors less then M", func(t *testing.T) {
		tasksCount := 20
		tasks := make([]Task, 0, tasksCount)
		errorsTasks := make([]Task, 0, tasksCount/2)
		tasksWithoutErrors := make([]Task, 0, tasksCount/2)

		var runTasksCount int32
		var sumTime time.Duration

		fillTasksListWithErrors(10, &errorsTasks, &runTasksCount)
		fillTasksList(10, &sumTime, &tasksWithoutErrors, &runTasksCount)
		tasks = append(tasks, errorsTasks...)
		tasks = append(tasks, tasksWithoutErrors...)
		workersCount := 10
		maxErrorsCount := 11

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}
