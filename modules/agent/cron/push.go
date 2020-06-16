package cron

import (
	"log"
	"time"

	"github.com/flytd/urlooker/dataobj"

	"github.com/flytd/urlooker/modules/agent/backend"
	"github.com/flytd/urlooker/modules/agent/g"
)

func Push() {
	for {
		checkResults := make([]*dataobj.CheckResult, 0)
		itemResults := g.CheckResultQueue.PopBack(500)
		if len(itemResults) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		for _, itemResult := range itemResults {
			checkResult := itemResult.(*dataobj.CheckResult)
			checkResults = append(checkResults, checkResult)
		}

		var resp string
		sendResultReq := dataobj.SendResultReq{
			Ip:           g.IP,
			CheckResults: checkResults,
		}
		err := backend.CallRpc("Web.SendResult", sendResultReq, &resp)
		if err != nil {
			log.Println("error:", err)
		}

		if g.Config.Debug {
			log.Println("<=", resp)
		}
	}
	return
}
