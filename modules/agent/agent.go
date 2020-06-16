package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/flytd/urlooker/modules/agent/backend"
	"github.com/flytd/urlooker/modules/agent/cron"
	"github.com/flytd/urlooker/modules/agent/g"
	"github.com/flytd/urlooker/modules/agent/utils"

	"github.com/toolkits/file"
)

func prepare() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func init() {
	prepare()

	cfg := flag.String("c", "", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	handleVersion(*version)
	handleHelp(*help)
	handleConfig(*cfg)

	fmt.Println("g.Config.Web.Addrs: ", g.Config)

	backend.InitClients(g.Config.Web.Addrs)

	g.Init()
}

func main() {
	ips, err := utils.IntranetIP()
	if len(ips) > 0 {
		g.IP = ips[0]
	} else {
		log.Println("get local ip err: ", err)
	}

	go cron.Push()
	cron.StartCheck()
}

func handleVersion(displayVersion bool) {
	if displayVersion {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
}

func handleHelp(displayHelp bool) {
	if displayHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func handleConfig(configFile string) {
	if configFile == "" {
		configFile = "configs/agent.yml"
	}

	if file.IsExist("configs/agent.local.yml") {
		configFile = "configs/agent.local.yml"
	}

	err := g.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
