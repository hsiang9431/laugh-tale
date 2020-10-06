package task

import (
    "os"
    "os/signal"
    "syscall"
	"github.com/pkg/errors"
)

func (tr *AsyncTaskRunner) Init(name string) error {
	if t, ok := tr.Tasks[name]; ok && t != nil {
		if t.Done() {
			return nil
		}
		return errors.Wrap(t.Init(), "Failed to init task: "+name)
	}
	return errors.New("No task named "+name)
}

func (tr *AsyncTaskRunner) InitAll() error {
	for n, _ := range tr.Tasks {
		if err := tr.Init(n); err != nil {
			return err
		}
	}
	return nil
}

func (tr *AsyncTaskRunner) Run(name string) error {
	if t, ok := tr.Tasks[name]; ok {
		tr.run(name, t)
		return nil
	}
	return errors.New("No task named "+name)
}

func (tr *AsyncTaskRunner) RunAll() {
	for n, t := range tr.Tasks {
		tr.run(n, t)
	}
}

func (tr *AsyncTaskRunner) WaitAll() []error {
	tr.wg.Wait()
	return tr.retErrs
}

func (tr *AsyncTaskRunner) run(name string, t Task) {
	if t == nil {
		return
	}
	if !t.Done() {
		var runErr error
		var cleanUpErr error
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		tr.wg.Add(1)
		go func() {
			go func() {
				runErr = errors.Wrap(t.Run(), "Failed to run task: "+name)
				close(c)
			}()
			<-c
			cleanUpErr = errors.Wrap(t.CleanUp(), "Failed to clean up task: "+name)
			
			if runErr != nil || cleanUpErr != nil {
				tr.retErrsLock.Lock()
				if runErr != nil {
					tr.retErrs = append(tr.retErrs, runErr)
				}
				if cleanUpErr != nil {
					tr.retErrs = append(tr.retErrs, cleanUpErr)
				}
				tr.retErrsLock.Unlock()
			}
			tr.wg.Done()
		}()
	}
}