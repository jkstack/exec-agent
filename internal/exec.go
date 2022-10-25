package internal

import (
	"context"
	"exec/internal/exec"
	"os"
	"strings"
	"time"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
)

func contain(cmd string, args []string, rm string) bool {
	if strings.Contains(cmd, rm) {
		return true
	}
	for _, arg := range args {
		if strings.Contains(arg, rm) {
			return true
		}
	}
	return false
}

func (agent *Agent) Run(msg *anet.Msg) error {
	if len(msg.Exec.DeferRM) > 0 && contain(msg.Exec.Cmd, msg.Exec.Args, msg.Exec.DeferRM) {
		defer os.Remove(msg.Exec.DeferRM)
	}
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
