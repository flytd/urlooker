package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/http/errors"
	"github.com/710leo/urlooker/modules/web/http/param"
	"github.com/710leo/urlooker/modules/web/http/render"
	"github.com/710leo/urlooker/modules/web/model"
	"github.com/710leo/urlooker/modules/web/utils"
)

type Strategy struct {
	Idc        string   `json:"idc"`
	Method     string   `json:"method"`
	Urls       []string `json:"urls"`
	Enable     int      `json:"enable"`
	Ip         string   `json:"ip"`
	Keywords   string   `json:"keywords"`
	Timeout    int      `json:"timeout"`
	ExpectCode string   `json:"expect_code"`
	Note       string   `json:"note"`
	Data       string   `json:"data"`
	Endpoint   string   `json:"endpoint"`
	Header     string   `json:"header"`
	PostData   string   `json:"post_data"`
	Times      int      `json:"times"`
	MaxStep    int      `json:"max_step"`
	Teams      string   `json:"teams"`
	Tags       []string `json:"tag"`
}

type StrUrl struct {
	Url 	[]string	`json:"url"`
}

func AddStrategyGet(w http.ResponseWriter, r *http.Request) {
	render.Put(r, "Regions", g.Config.IDC)
	render.HTML(r, w, "strategy/create")
}

func AddStrategyPost(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))
	var msg string
	var err error
	var tagStr string
	urls := strings.Split(param.String(r, "url", ""), "\n")
	for _, url := range urls {
		err := utils.CheckUrl(url)
		if err != nil {
			errors.Panic(fmt.Sprintf("url:%s %s", url, err.Error()))
		}
	}

	tags := strings.Split(param.String(r, "tags", ""), "\n")
	if len(tags) > 0 && tags[0] != "" {
		for _, tag := range tags {
			strS := strings.Split(tag, "=")
			if len(strS) != 2 {
				errors.Panic("tag must be like aaa=bbb")
			}
		}
		tagStr = strings.Join(tags, ",")
	}

	for _, url := range urls {
		var s = model.Strategy{}
		s.Method = param.String(r, "method", "get")
		s.Creator = me.Name
		s.Enable = param.Int(r, "enable", 1)
		s.Url = url
		s.Idc = param.String(r, "idc", "default")
		s.ExpectCode = param.String(r, "expect_code", "200")
		s.Timeout = param.Int(r, "timeout", 3000)
		s.Header = param.String(r, "header", "")
		s.PostData = param.String(r, "post_data", "")
		s.MaxStep = param.Int(r, "max_step", 3)
		s.Teams = param.MustString(r, "teams")
		s.Times = param.Int(r, "times", 3)
		s.Note = param.String(r, "note", "")
		s.Keywords = param.String(r, "keywords", "")
		s.Data = param.String(r, "data", "")
		s.Endpoint = param.String(r, "endpoint", "")
		s.Tag = tagStr
		s.IP = param.String(r, "ip", "")

		_, err = s.Add()
		if err != nil {
			msg += fmt.Sprintf("strategy:%s failed, err:%s", url, err.Error())
		} else {
			msg += fmt.Sprintf("strategy:%s success :)", url)
		}
	}

	//errors.MaybePanic(err)
	if err != nil {
		errMsg := fmt.Sprintf("%s,err:%v", msg, err)
		errors.Panic(errMsg)
	}
	render.Data(w, msg)
}

func AddStrategyApiPost(w http.ResponseWriter, r *http.Request) {
	user := IsLogin(r)
	if user == nil || user.Name == "" {
		render.ErrorCode(w, errors.NewError("没有用户登录"))
		return
	}
	username := user.Name
	var msg string
	var err error
	var tagStr string

	result, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var data Strategy
	json.Unmarshal(result, &data)
	urls := data.Urls
	for _, url := range urls {
		err := utils.CheckUrl(url)
		if err != nil {
			render.ErrorCode(w, err)
			return
		}
	}
	tags := data.Tags
	if len(tags) > 0 && tags[0] != "" {
		for _, tag := range tags {
			strS := strings.Split(tag, "=")
			if len(strS) != 2 {
				render.ErrorCode(w, errors.NewError("tag 必须类似于 aaa=bbb"))
				return
			}
		}
		tagStr = strings.Join(tags, ",")
	}

	for _, url := range urls {
		var s = model.Strategy{}
		s.Method = param.ApiString(data.Method, "get")
		s.Creator = username
		s.Enable = data.Enable
		s.Url = url
		s.Idc = param.ApiString(data.Idc, "default")
		s.ExpectCode = param.ApiString(data.ExpectCode, "200")
		s.Timeout = param.ApiInt(data.Timeout, 3000)
		s.Header = param.ApiString(data.Header, "")
		s.PostData = param.ApiString(data.PostData, "")
		s.MaxStep = param.ApiInt(data.MaxStep, 3)
		s.Teams = data.Teams
		s.Times = param.ApiInt(data.Times, 3)
		s.Note = param.ApiString(data.Note, "")
		s.Keywords = param.ApiString(data.Keywords, "")
		s.Data = param.ApiString(data.Data, "")
		s.Endpoint = param.ApiString(data.Endpoint, "")
		s.Tag = tagStr
		s.IP = param.ApiString(data.Ip, "")

		_, err = s.Add()
		if err != nil {
			msg += fmt.Sprintf("添加:%s 失败, 错误:%s", url, err.Error())
		} else {
			msg += fmt.Sprintf("添加%s成功 ", url)
		}
	}

	if err != nil {
		render.ErrorCode(w, errors.NewError(msg))
		return
	}
	render.Data(w, msg)
}

func GetStrategyById(w http.ResponseWriter, r *http.Request) {
	strategy := StraRequired(r)
	render.Data(w, strategy)
}

func UpdateStrategyGet(w http.ResponseWriter, r *http.Request) {
	s := StraRequired(r)
	render.Put(r, "Id", s.Id)
	render.Put(r, "Regions", g.Config.IDC)
	render.HTML(r, w, "strategy/edit")
}

func UpdateStrategy(w http.ResponseWriter, r *http.Request) {
	s := StraRequired(r)
	me := MeRequired(LoginRequired(w, r))
	var tagStr string

	username := me.Name
	if s.Creator != username && !IsAdmin(username) {
		errors.Panic("没有权限")
	}

	url := param.String(r, "url", "")
	err := utils.CheckUrl(url)
	if err != nil {
		errors.Panic(fmt.Sprintf("url:%s %s", url, err.Error()))
	}

	tags := strings.Split(param.String(r, "tags", ""), "\n")
	if len(tags) > 0 && tags[0] != "" {
		log.Println("len:", len(tags))
		for _, tag := range tags {
			strs := strings.Split(tag, "=")
			if len(strs) != 2 {
				errors.Panic("tag must be like aaa=bbb")
			} else if strs[0] == "" || strs[1] == "" {
				errors.Panic("tag must be like aaa=bbb")
			}
		}
		tagStr = strings.Join(tags, ",")
	}

	s.Creator = username
	s.Url = url
	s.Idc = param.String(r, "idc", "default")
	s.Method = param.String(r, "method", "get")
	s.Enable = param.Int(r, "enable", 1)
	s.ExpectCode = param.String(r, "expect_code", "200")
	s.Timeout = param.Int(r, "timeout", 3000)
	s.Header = param.String(r, "header", "")
	s.PostData = param.String(r, "post_data", "")
	s.MaxStep = param.Int(r, "max_step", 3)
	s.Teams = param.String(r, "teams", "")
	s.Times = param.Int(r, "times", 3)
	s.Note = param.String(r, "note", "")
	s.Keywords = param.String(r, "keywords", "")
	s.Data = param.String(r, "data", "")
	s.Endpoint = param.String(r, "endpoint", "")
	s.IP = param.String(r, "ip", "")
	s.Tag = tagStr

	err = s.Update()
	errors.MaybePanic(err)
	render.Data(w, "ok")
}

func DeleteStrategy(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))
	strategy := StraRequired(r)
	//teams := strings.Split(strategy.Teams, ",")

	username := me.Name
	if strategy.Creator != username && !IsAdmin(username) {
		errors.Panic("没有权限")
	}

	err := strategy.Delete()
	errors.MaybePanic(err)
	render.Data(w, "ok")
}

func DeleteStrategyApi(w http.ResponseWriter, r *http.Request) {
	user := IsLogin(r)
	if user == nil || user.Name == "" {
		render.ErrorCode(w, errors.NewError("没有用户登录"))
		return
	}
	var msg string
	var err error
	var data StrUrl
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &data)
	for _, url := range data.Url {
		strategy, err := model.GetStrategyByUrl(url)
		username := user.Name
		if strategy.Creator != username && !IsAdmin(username) {
			render.ErrorCode(w, errors.NewError("没有权限"))
			return
		}
		err = strategy.Delete()
		if err != nil {
			msg += fmt.Sprintf("删除:%s 失败, 错误:%s", url, err.Error())
		} else {
			msg += fmt.Sprintf("删除%s成功 ", url)
		}
	}
	if err != nil {
		render.ErrorCode(w, errors.NewError(msg))
	}
	render.Data(w, msg)
}

func GetTeamsOfStrategy(w http.ResponseWriter, r *http.Request) {
	MeRequired(LoginRequired(w, r))
	stra := StraRequired(r)
	teams, err := model.GetTeamsByIds(stra.Teams)
	errors.MaybePanic(err)
	render.Data(w, map[string]interface{}{"teams": teams})
}
