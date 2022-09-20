package internal

import (
	"os"

	"github.com/jkstack/jkframe/conf/kvconf"
	"github.com/jkstack/jkframe/utils"
	"github.com/jkstack/libagent/conf"
)

type Configure conf.Configure

func load(dir string) *Configure {
	f, err := os.Open(dir)
	utils.Assert(err)
	defer f.Close()
	var ret Configure
	utils.Assert(kvconf.NewDecoder(f).Decode(&ret))
	return &ret
}

func (agent *Agent) ConfDir() string {
	return agent.cfgDir
}

func (agent *Agent) Configure() *conf.Configure {
	return (*conf.Configure)(agent.cfg)
}

func (agent *Agent) OnRewriteConfigure() error {
	f, err := os.Create(agent.cfgDir + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	err = kvconf.NewEncoder(f).Encode(agent.cfg)
	if err != nil {
		return err
	}
	f.Close()
	os.Rename(agent.cfgDir+".tmp", agent.cfgDir)
	return nil
}
