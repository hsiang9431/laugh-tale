package task

import "sync"

type Task interface {
	Init() error
	Run() error
	CleanUp() error
	Done() bool
}

type TaskRunner struct {
	Tasks map[string]Task
}

type AsyncTaskRunner struct {
	Tasks map[string]Task

	retErrs     []error
	retErrsLock sync.Mutex
	wg          sync.WaitGroup
}
