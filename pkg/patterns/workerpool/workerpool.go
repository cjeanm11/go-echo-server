package taskmanager

import (
	"sync"
)

type Item interface{}

type TaskCommand interface {
	invoke() interface{}
}

type Task[T any] struct {
	execute func() T
}

func NewTask[T any](execute func() T) Task[T] {
	return Task[T]{execute: execute}
}

func (t Task[T]) invoke() interface{} {
	return t.execute()
}

func ProcessTasks(numWorkers int, tasks []func() interface{}) []interface{} {
	var wg sync.WaitGroup
	taskChannel := make(chan func() interface{}, len(tasks))
	resultChannel := make(chan interface{}, len(tasks))
	wg.Add(numWorkers)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		go func() {
			for task := range taskChannel {
				result := task()
				resultChannel <- result
			}
			wg.Done()
		}()
	}

	// Submit tasks
	for _, task := range tasks {
		taskChannel <- task
	}
	close(taskChannel)

	// Wait for all workers to finish
	wg.Wait()

	// Close the result channel after all results are received
	close(resultChannel)

	// Collect results
	var results []interface{}
	for result := range resultChannel {
		results = append(results, result)
	}

	return results
}
