package internal

import (
	"exec/internal/exec"
	"sync"

	"github.com/jkstack/anet"
)

// AgentName agent name
var AgentName string

// Agent agent object
type Agent struct {
	sync.RWMutex
	cfgDir  string
	cfg     *Configure
	version string
	chWrite chan *anet.Msg
	// runtime
	tasks map[int]*exec.Task
}

// New create agent object
func New(dir, version string) *Agent {
	return &Agent{
		cfgDir:  dir,
		cfg:     load(dir),
		version: version,
		chWrite: make(chan *anet.Msg),
		tasks:   make(map[int]*exec.Task),
	}
}

// AgentName get agent name
func (agent *Agent) AgentName() string {
	return "exec-agent"
}

// Version geet agent version
func (agent *Agent) Version() string {
	return agent.version
}
