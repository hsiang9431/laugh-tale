package task

import "github.com/pkg/errors"


func (tr *TaskRunner) Init(name string) error {
	if t, ok := tr.Tasks[name]; ok && t != nil {
		if t.Done() {
			return nil
		}
		return errors.Wrap(t.Init(), "Failed to init task: "+name)
	}
	return errors.New("No task named "+name)
}

func (tr *TaskRunner) InitAll() error {
	for n, _ := range tr.Tasks {
		if err := tr.Init(n); err != nil {
			return err
		}
	}
	return nil
}

func (tr *TaskRunner) Run(name string) error {
	if t, ok := tr.Tasks[name]; ok && t != nil {
		if t.Done() {
			return nil
		}
		if err := t.Run(); err != nil {
			return errors.Wrap(err, "Failed to run task: "+name)
		}
		if err := t.CleanUp(); err != nil {
			return errors.Wrap(err, "Failed to clean up task: "+name)
		}
		return nil
	}
	return errors.New("No task named "+name)
}

func (tr *TaskRunner) RunAll() error {
	for n, _ := range tr.Tasks {
		if err := tr.Run(n); err != nil {
			return err
		}
	}
	return nil
}
