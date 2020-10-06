package task

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type testTask struct {
	Name         string
	HasInit      bool
	HasRun       bool
	HasCleanedUp bool
	Runtime      int

	TestingError bool

	T *testing.T
}

func (tt *testTask) Done() bool {
	return tt.HasRun && tt.HasCleanedUp
}

func (tt *testTask) Init() error {
	tt.T.Log("call init on " + tt.Name)
	if tt.TestingError {
		tt.HasInit = true
		return errors.New("test error: run " + tt.Name)
	}
	tt.HasInit = true
	tt.Runtime = rand.Intn(200)
	return nil
}

func (tt *testTask) Run() error {
	tt.T.Log("call run on " + tt.Name)
	if !tt.HasInit {
		return errors.New("run failed: task not init " + tt.Name)
	}
	if tt.TestingError {
		tt.HasRun = true
		return errors.New("test error: run " + tt.Name)
	}
	time.Sleep(time.Duration(tt.Runtime+200) * time.Millisecond)
	tt.HasRun = true
	tt.T.Log("run done on " + tt.Name)
	return nil
}

func (tt *testTask) CleanUp() error {
	tt.T.Log("call clean up on " + tt.Name)
	if !tt.HasRun {
		return errors.New("clean up failed: task not run " + tt.Name)
	}
	if tt.TestingError {
		return errors.New("test error: clean up " + tt.Name)
	}
	time.Sleep(time.Duration(tt.Runtime+100) * time.Millisecond)
	tt.HasCleanedUp = true
	tt.T.Log("clean up done on " + tt.Name)
	return nil
}

func getTasks(cnt int, testingErr bool, t *testing.T) map[string]Task {
	ret := make(map[string]Task)
	for i := 0; i < cnt; i++ {
		name := fmt.Sprintf("test-task-%02d", i+1)
		ret[name] = &testTask{Name: name, HasInit: testingErr, TestingError: testingErr, T: t}
	}
	return ret
}

func TestSync(t *testing.T) {
	tasks := getTasks(5, false, t)
	tr := TaskRunner{
		Tasks: tasks,
	}
	if err := tr.InitAll(); err != nil {
		t.Error("TestSync failed: init error")
	}
	if err := tr.RunAll(); err != nil {
		t.Error("TestSync failed: run error")
	}
	for _, tsk := range tasks {
		if !tsk.Done() {
			errMsg := fmt.Sprintf("TestSync failed: task %v not done", tsk)
			t.Error(errMsg)
		}
	}
}

func TestAsync(t *testing.T) {
	tasks := getTasks(5, false, t)
	tr := AsyncTaskRunner{
		Tasks: tasks,
	}
	if err := tr.InitAll(); err != nil {
		t.Error("TestAsync failed: init error")
	}
	tr.RunAll()
	errs := tr.WaitAll()
	if len(errs) != 0 {
		t.Error("TestSync failed: run error")
	}
	for _, tsk := range tasks {
		if !tsk.Done() {
			errMsg := fmt.Sprintf("TestSync failed: task %v not done", tsk)
			t.Error(errMsg)
		}
	}
}

func TestRunError(t *testing.T) {
	tasks1 := getTasks(1, true, t)
	tr := TaskRunner{
		Tasks: tasks1,
	}
	if err := tr.InitAll(); err == nil {
		t.Error("sync init error not captured")
	}
	if err := tr.RunAll(); err == nil {
		t.Error("sync run error not captured")
	}
	tasks2 := getTasks(5, true, t)
	atr := AsyncTaskRunner{
		Tasks: tasks2,
	}
	if err := atr.InitAll(); err == nil {
		t.Error("async init error not captured")
	}
	atr.RunAll()
	errs := atr.WaitAll()
	if len(errs) == 0 {
		t.Error("async run error not captured")
	}
	for _, err := range errs {
		t.Log(err.Error())
	}
}
