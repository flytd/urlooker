package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/toolkits/str"

	"github.com/flytd/urlooker/modules/web/g"
	"github.com/flytd/urlooker/modules/web/http/cookie"
	"github.com/flytd/urlooker/modules/web/http/errors"
	"github.com/flytd/urlooker/modules/web/model"
)

func StraRequired(r *http.Request) *model.Strategy {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	errors.MaybePanic(err)

	obj, err := model.GetStrategyById(id)
	errors.MaybePanic(err)
	if obj == nil {
		panic(errors.BadRequestError("no such item"))
	}
	return obj
}

func BindJson(r *http.Request, obj interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("Empty request body")
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, obj)
	if err != nil {
		return fmt.Errorf("unmarshal body %s err:%v", string(body), err)
	}
	return err
}

func IdcRequired(r *http.Request) string {
	vars := mux.Vars(r)
	idc := vars["idc"]

	if str.HasDangerousCharacters(idc) {
		errors.Panic("idc不合法")
	}

	return idc
}

func LoginRequired(w http.ResponseWriter, r *http.Request) (int64, string) {
	userId, username, found := cookie.ReadUser(r)
	if !found {
		panic(errors.NotLoginError())
	}

	return userId, username
}

func IsLogin (r *http.Request) *model.User {
	userId, _, found := cookie.ReadUser(r)
	user, err := model.GetUserById(userId)
	if err != nil || !found || user == nil {
		return nil
	}
	return user
}


func AdminRequired(id int64, name string) {
	user, err := model.GetUserById(id)
	if err != nil {
		panic(errors.InternalServerError(err.Error()))
	}

	if user == nil {
		panic(errors.NotLoginError())
	}

	for _, admin := range g.Config.Admins {
		if user.Name == admin {
			return
		}
	}

	panic(errors.NotLoginError())
	return
}

func MeRequired(id int64, name string) *model.User {
	user, err := model.GetUserById(id)
	if err != nil {
		panic(errors.InternalServerError(err.Error()))
	}

	if user == nil {
		panic(errors.NotLoginError())
	}

	return user
}

func TeamRequired(r *http.Request) *model.Team {
	vars := mux.Vars(r)
	tid, err := strconv.ParseInt(vars["tid"], 10, 64)
	errors.MaybePanic(err)

	team, err := model.GetTeamById(tid)
	errors.MaybePanic(err)
	if team == nil {
		panic(errors.BadRequestError("no such team"))
	}

	return team
}

func UserMustBeMemberOfTeam(uid, tid int64) {
	is, err := model.IsMemberOfTeam(uid, tid)
	errors.MaybePanic(err)
	if is {
		return
	}

	team, err := model.GetTeamById(tid)
	errors.MaybePanic(err)
	if team != nil && team.Creator == uid {
		return
	}

	panic(errors.BadRequestError("用户不是团队的成员"))
}

func IsAdmin(username string) bool {
	for _, admin := range g.Config.Admins {
		if username == admin {
			return true
		}
	}
	return false
}
