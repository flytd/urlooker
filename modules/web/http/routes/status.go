package routes

import (
	"net/http"

	"github.com/flytd/urlooker/modules/web/g"
	"github.com/flytd/urlooker/modules/web/http/errors"
	"github.com/flytd/urlooker/modules/web/http/render"
	"github.com/flytd/urlooker/modules/web/utils"
)

func GetLog(w http.ResponseWriter, r *http.Request) {

	AdminRequired(LoginRequired(w, r))
	appLog, err := utils.ReadLastLine("var/app.log")
	errors.MaybePanic(err)

	render.Put(r, "Log", appLog)

	render.HTML(r, w, "status/log")
}

func Version(w http.ResponseWriter, r *http.Request) {
	render.Data(w, g.VERSION)
}

func Health(w http.ResponseWriter, r *http.Request) {
	render.Data(w, "ok")
}
