package http

import (
	"github.com/flytd/urlooker/modules/web/g"
	"github.com/flytd/urlooker/modules/web/http/middleware"
	"github.com/flytd/urlooker/modules/web/http/render"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func Start() {
	render.Init()

	r := mux.NewRouter().StrictSlash(false)

	ConfigRouter(r)
	n := negroni.New()
	n.Use(middleware.NewLogger())
	n.Use(middleware.NewRecovery())
	n.UseHandler(r)
	n.Run(g.Config.Http.Listen)
}
