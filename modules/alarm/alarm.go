package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/flytd/urlooker/modules/alarm/backend"
	"github.com/flytd/urlooker/modules/alarm/cron"
	"github.com/flytd/urlooker/modules/alarm/g"
	"github.com/flytd/urlooker/modules/alarm/judge"
	"github.com/flytd/urlooker/modules/alarm/receiver"
	"github.com/flytd/urlooker/modules/alarm/sender"

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

	judge.InitHistoryBigMap()
	sender.Init()
	backend.InitClients(g.Config.Web.Addrs)
}

func main() {
	go cron.SyncStrategies()
	receiver.Start()
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
		configFile = "configs/alarm.yml"
	}

	if file.IsExist("configs/alarm.local.yml") {
		configFile = "configs/alarm.local.yml"
	}

	err := g.Parse(configFile)
	if err != nil {
		log.Fatalln(err)
	}
}
