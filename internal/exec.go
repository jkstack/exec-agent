package internal

import (
	"context"
	"exec/internal/exec"
	"time"

	"github.com/jkstack/anet"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go task.Response(agent.chWrite, cancel)
	go task.Wait(ctx, agent.chWrite)
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
