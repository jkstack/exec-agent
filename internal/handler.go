package internal

import (
	"context"

	"github.com/jkstack/anet"
)

func (agent *Agent) OnConnect() {
}

func (agent *Agent) OnDisconnect() {
}

func (agent *Agent) OnReportMonitor() {
}

func (agent *Agent) OnMessage(msg *anet.Msg) error {
	return nil
}

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
