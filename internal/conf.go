package internal

import (
	"os"

	"github.com/jkstack/jkframe/conf/kvconf"
	"github.com/jkstack/jkframe/utils"
	"github.com/jkstack/libagent/conf"
)

// Configure configure object
type Configure struct {
	Basic conf.Configure `kv:"basic"`
}

func load(dir string) *Configure {
	f, err := os.Open(dir)
	utils.Assert(err)
	defer f.Close()
	var ret Configure
	utils.Assert(kvconf.NewDecoder(f).Decode(&ret))
	return &ret
}

// ConfDir get configure file dir
func (agent *Agent) ConfDir() string {
	return agent.cfgDir
}

// Configure get basic configure
func (agent *Agent) Configure() *conf.Configure {
	return (*conf.Configure)(&agent.cfg.Basic)
}

// OnRewriteConfigure rewrite configure file
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
