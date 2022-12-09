package internal

import (
	"context"

	"github.com/jkstack/anet"
)

// OnConnect on connect callback
func (agent *Agent) OnConnect() {
}

// OnDisconnect on disconnect callback
func (agent *Agent) OnDisconnect() {
}

// OnReportMonitor on report monitor callback
func (agent *Agent) OnReportMonitor() {
}

// OnMessage on receive message callback
func (agent *Agent) OnMessage(msg *anet.Msg) error {
	switch msg.Type {
	case anet.TypeExec:
		return agent.Run(msg)
	case anet.TypeExecKill:
		return agent.Kill(msg.ExecKill.Pid)
	case anet.TypeLsReq:
		return agent.Ls(msg.TaskID, msg.LSReq.Dir)
	case anet.TypeDownloadReq:
		return agent.Download(msg.TaskID, msg.DownloadReq.Dir)
	case anet.TypeUpload:
		return agent.Upload(msg.TaskID, msg)
	}
	return nil
}

// LoopWrite loop write
func (agent *Agent) LoopWrite(ctx context.Context, ch chan *anet.Msg) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-agent.chWrite:
			ch <- msg
		}
	}
}
