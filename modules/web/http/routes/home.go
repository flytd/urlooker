package routes

import (
	"github.com/toolkits/str"
	"github.com/toolkits/web"
	"net/http"

	"github.com/710leo/urlooker/modules/web/http/errors"
	"github.com/710leo/urlooker/modules/web/http/param"
	"github.com/710leo/urlooker/modules/web/http/render"
	"github.com/710leo/urlooker/modules/web/model"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))
	username := me.Name
	mine := param.Int(r, "mine", 1)
	query := param.String(r, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
	}
	if IsAdmin(username) {
		mine = 0
	}

	limit := param.Int(r, "limit", 10)
	total, err := model.GetAllStrategyCount(mine, query, username)
	errors.MaybePanic(err)
	pager := web.NewPaginator(r, limit, total)

	strategies, err := model.GetAllStrategy(mine, limit, pager.Offset(), query, username)

	errors.MaybePanic(err)
	render.Put(r, "Strategies", strategies)
	render.Put(r, "Pager", pager)
	render.Put(r, "Mine", mine)
	render.Put(r, "Query", query)
	render.HTML(r, w, "home/index")
}

func HomeApiIndex(w http.ResponseWriter, r *http.Request) {
	user := IsLogin(r)
	if user == nil || user.Name == "" {
		render.ErrorCode(w,errors.NewError("没有用户登录"))
		return
	}
	username := user.Name
	mine := param.Int(r, "mine", 1)
	query := param.String(r, "q", "")
	if str.HasDangerousCharacters(query) {
		errors.Panic("查询字符不合法")
		render.ErrorCode(w,errors.NewError("查询字符不合法"))
		return
	}
	if IsAdmin(username) {
		mine = 0
	}

	limit := param.Int(r, "limit", 10)
	total, err := model.GetAllStrategyCount(mine, query, username)
	if err != nil {
		errors.MaybePanic(err)
		render.ErrorCode(w, err)
		return
	}
	pager := web.NewPaginator(r, limit, total)
	strategies, err := model.GetAllStrategy(mine, limit, pager.Offset(), query, username)
	if err != nil {
		errors.MaybePanic(err)
		render.ErrorCode(w, err)
		return
	}
	data := make(map[string]interface{})
	data["strategies"] = strategies
	data["page"] = pager.Page()
	data["pagenums"] = pager.PageNums()
	msg := "获取成功"
	render.Data(w, data, msg)
}