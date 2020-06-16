package routes

import (
	"net/http"

	"github.com/flytd/urlooker/modules/web/g"
	"github.com/flytd/urlooker/modules/web/http/render"
)

func GetDetectItem(w http.ResponseWriter, r *http.Request) {
	detectItem, _ := g.DetectedItemMap.Get(IdcRequired(r))
	render.Data(w, detectItem)
}
