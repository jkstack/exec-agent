package internal

import (
	"context"
	"exec/internal/exec"
	"time"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
)

func (agent *Agent) Run(msg *anet.Msg) error {
	task := exec.NewTask(msg.TaskID)
	err := task.Prepare(msg.Exec)
	if err != nil {
		agent.execError(task, err.Error())
		return err
	}
	err = task.Start(time.Duration(msg.Exec.Timeout) * time.Second)
	if err != nil {
		agent.execError(task, err.Error())
		return err
	}
	agent.execOK(task)

	agent.Lock()
	agent.tasks[task.Pid()] = task
	agent.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go task.Response(agent.chWrite, cancel)
	go func() {
		task.Wait(ctx, agent.chWrite)
		task.Close()

		agent.Lock()
		delete(agent.tasks, task.Pid())
		agent.Unlock()
	}()
	return nil
}

func (agent *Agent) execError(task *exec.Task, msg string) {
	var m anet.Msg
	m.TaskID = task.ID
	m.Type = anet.TypeExecd
	m.Execd = &anet.ExecdPayload{
		OK:   false,
		Msg:  msg,
		Pid:  task.Pid(),
		Time: task.Begin,
	}
	agent.chWrite <- &m
}

func (agent *Agent) execOK(task *exec.Task) {
	var msg anet.Msg
	msg.TaskID = task.ID
	msg.Type = anet.TypeExecd
	msg.Execd = &anet.ExecdPayload{
		OK:   true,
		Pid:  task.Pid(),
		Time: task.Begin,
	}
	agent.chWrite <- &msg
}

func (agent *Agent) Kill(pid int) error {
	agent.RLock()
	task := agent.tasks[pid]
	agent.RUnlock()
	if task == nil {
		logging.Warning("process of %d not found", pid)
		return nil
	}
	logging.Info("kill for task [%s], pid=%d", task.ID, pid)
	task.Close()
	return nil
}
