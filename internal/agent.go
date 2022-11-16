package internal

import (
	"exec/internal/exec"
	"sync"

	"github.com/jkstack/anet"
)

var AgentName string

type Agent struct {
	sync.RWMutex
	cfgDir  string
	cfg     *Configure
	version string
	chWrite chan *anet.Msg
	// runtime
	tasks map[int]*exec.Task
}

func New(dir, version string) *Agent {
	return &Agent{
		cfgDir:  dir,
		cfg:     load(dir),
		version: version,
		chWrite: make(chan *anet.Msg),
		tasks:   make(map[int]*exec.Task),
	}
}

func (agent *Agent) AgentName() string {
	return "exec-agent"
}

func (agent *Agent) Version() string {
	return agent.version
}
