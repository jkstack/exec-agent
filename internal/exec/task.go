package exec

import (
	"bytes"
	"context"
	"encoding/base64"
	"exec/internal/utils"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Task task object
type Task struct {
	ID        string
	Begin     time.Time
	cmd       *exec.Cmd
	stdin     string
	pid       int
	tty       io.Reader
	closeFunc func()
	done      chan struct{}
}

// NewTask create task object
func NewTask(id string) *Task {
	return &Task{
		ID:    id,
		Begin: time.Now(),
		done:  make(chan struct{}),
	}
}

// Close close task
func (t *Task) Close() {
	select {
	case t.done <- struct{}{}:
	default:
	}
	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}
	if t.closeFunc != nil {
		t.closeFunc()
	}
}

// Prepare prepare task
func (t *Task) Prepare(req *anet.ExecPayload) error {
	if req.Timeout <= 0 {
		logging.Warning("reset timeout to 60s for task %s", t.ID)
		req.Timeout = 60
	}
	var err error
	t.cmd, t.stdin, err = buildCommand(req)
	if err != nil {
		return err
	}
	logging.Info("build command for task: %s\n  => %s", t.ID, t.cmd.String())
	if len(req.WorkDir) > 0 {
		t.cmd.Dir = req.WorkDir
		logging.Info("set working directory to [%s] for task: %s", req.WorkDir, t.ID)
	}
	t.cmd.Env = os.Environ()
	if len(req.Env) > 0 {
		t.cmd.Env = append(t.cmd.Env, req.Env...)
		env := make([]string, len(req.Env))
		for i, e := range req.Env {
			env[i] = "  ==> " + e
		}
		logging.Info("add environment for task: %s\n%s", t.ID, strings.Join(env, "\n"))
	}
	return nil
}

// Start start task
func (t *Task) Start(timeout time.Duration) error {
	var err error
	t.tty, t.closeFunc, err = createPty(t.cmd, t.stdin)
	if err != nil {
		return err
	}
	if t.cmd.Process != nil {
		t.pid = t.cmd.Process.Pid
		logging.Info("start process for task [%s] success, pid=%d", t.ID, t.pid)
	}
	go func() {
		select {
		case <-time.After(timeout):
			logging.Info("task %s timeouted after %d seconds", t.ID, int(timeout.Seconds()))
			if t.cmd != nil {
				t.cmd.Process.Kill()
			}
		case <-t.done:
			return
		}
	}()
	return nil
}

// Wait wait task
func (t *Task) Wait(ctx context.Context, ch chan *anet.Msg) {
	if t.cmd == nil {
		return
	}
	var code int
	if err := t.cmd.Wait(); err != nil {
		if ex, ok := err.(*exec.ExitError); ok {
			code = ex.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			logging.Error("is not ExitError: %T", err)
		}
	}
	logging.Info("task %s run done, code=%d", t.ID, code)
	<-ctx.Done()
	var msg anet.Msg
	msg.TaskID = t.ID
	msg.Type = anet.TypeExecDone
	msg.ExecDone = &anet.ExecDone{
		Code: code,
		Time: time.Now(),
	}
	ch <- &msg
	select {
	case t.done <- struct{}{}:
	default:
	}
}

// Response response task
func (t *Task) Response(ch chan *anet.Msg, cancel context.CancelFunc) {
	defer cancel()
	buf := make([]byte, 64*1024)
	for {
		n, err := t.tty.Read(buf)
		if err != nil {
			return
		}
		send := buf[:n]
		if !utils.IsUtf8(send) {
			r := transform.NewReader(bytes.NewReader(send), simplifiedchinese.GBK.NewDecoder())
			var err error
			data, err := io.ReadAll(r)
			if err != nil {
				return
			}
			send = data
		}
		var msg anet.Msg
		msg.TaskID = t.ID
		msg.Type = anet.TypeExecData
		msg.ExecData = &anet.ExecData{
			Data: base64.StdEncoding.EncodeToString(send),
		}
		ch <- &msg
		logging.Debug("task %s sent %s data", t.ID, humanize.Bytes(uint64(len(send))))
	}
}

// Pid get pid
func (t *Task) Pid() int {
	return t.pid
}
