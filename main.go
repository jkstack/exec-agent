package main

import (
	"exec/internal"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	rt "runtime"

	"github.com/jkstack/jkframe/utils"
	agent "github.com/jkstack/libagent"
)

const agentName = "exec-agent"

var (
	version      string = "0.0.0"
	gitBranch    string = "<branch>"
	gitHash      string = "<hash>"
	gitReversion string = "0"
	buildTime    string = "0000-00-00 00:00:00"
)

func showVersion() {
	fmt.Printf("version: %s\ncode version: %s.%s.%s\nbuild time: %s\ngo version: %s\n",
		version,
		gitBranch, gitHash, gitReversion,
		buildTime,
		rt.Version())
}

func main() {
	cf := flag.String("conf", "", "config file dir")
	ver := flag.Bool("version", false, "show version info")
	act := flag.String("action", "", "install, uninstall")
	flag.Parse()

	if *ver {
		showVersion()
		return
	}

	switch *act {
	case "install":
		if len(*cf) == 0 {
			fmt.Println("missing -conf argument")
			os.Exit(1)
		}

		dir, err := filepath.Abs(*cf)
		utils.Assert(err)

		dummy := agent.NewDummyApp(agentName, dir)

		err = agent.RegisterService(dummy)
		if err != nil {
			fmt.Printf("can not register service: %v\n", err)
			return
		}
		fmt.Println("register service success")
	case "uninstall":
		err := agent.UnregisterService(agent.NewDummyApp(agentName, ""))
		if err != nil {
			fmt.Printf("can not unregister service: %v\n", err)
			return
		}
		fmt.Println("unregister service success")
	case "start":
		err := agent.Start(agent.NewDummyApp(agentName, ""))
		if err != nil {
			fmt.Printf("can not start service: %v\n", err)
			return
		}
		fmt.Println("start service success")
	case "stop":
		err := agent.Stop(agent.NewDummyApp(agentName, ""))
		if err != nil {
			fmt.Printf("can not stop service: %v\n", err)
			return
		}
		fmt.Println("stop service success")
	default:
		if len(*cf) == 0 {
			fmt.Println("missing -conf argument")
			os.Exit(1)
		}

		dir, err := filepath.Abs(*cf)
		utils.Assert(err)

		agent.Run(internal.New(dir, version))
	}
}
