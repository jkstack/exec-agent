package internal

import "github.com/jkstack/anet"

type Agent struct {
	cfgDir  string
	cfg     *Configure
	version string
	chWrite chan *anet.Msg
}

func New(dir, version string) *Agent {
	return &Agent{
		cfgDir:  dir,
		cfg:     load(dir),
		version: version,
		chWrite: make(chan *anet.Msg),
	}
}

func (agent *Agent) AgentName() string {
	return "exec-agent"
}

func (agent *Agent) Version() string {
	return agent.version
}
