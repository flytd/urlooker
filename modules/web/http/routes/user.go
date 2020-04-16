package routes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/toolkits/str"

	"github.com/710leo/urlooker/modules/web/g"
	"github.com/710leo/urlooker/modules/web/http/cookie"
	"github.com/710leo/urlooker/modules/web/http/errors"
	"github.com/710leo/urlooker/modules/web/http/param"
	"github.com/710leo/urlooker/modules/web/http/render"
	"github.com/710leo/urlooker/modules/web/model"
	"github.com/710leo/urlooker/modules/web/utils"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if g.Config.Ldap.Enabled || !g.Config.Register {
		errors.Panic("注册已关闭")
	}
	username := param.MustString(r, "username")
	password := param.MustString(r, "password")
	repeat := param.MustString(r, "repeat")

	if password != repeat {
		errors.Panic("两次输入的密码不一致")
	}

	if str.HasDangerousCharacters(username) {
		errors.Panic("用户名不合法，请不要使用非法字符")
	}

	userid, err := model.UserRegister(username, utils.EncryptPassword(g.Config.Salt, password))
	errors.MaybePanic(err)

	render.Data(w, cookie.WriteUser(w, userid, username))
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	render.Put(r, "Title", "register")
	render.Put(r, "callback", param.String(r, "callback", "/"))
	render.HTML(r, w, "auth/register")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	errors.MaybePanic(cookie.RemoveUser(w))
	http.Redirect(w, r, "/", 302)
}

func ApiLogout(w http.ResponseWriter, r *http.Request) {
	_, name, _ :=cookie.ReadUser(r)
	if name == "" {
		render.ErrorCode(w, errors.NewError( "没用用户登录"))
		return
	}
	if err := cookie.RemoveUser(w); err != nil {
		errors.MaybePanic(err)
		render.ErrorCode(w, err)
		return
	}
	render.Data(w, "", "登出成功")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	render.Put(r, "Title", "login")
	render.Put(r, "callback", param.String(r, "callback", "/"))
	render.HTML(r, w, "auth/login")
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := param.MustString(r, "username")
	password := param.MustString(r, "password")

	if str.HasDangerousCharacters(username) {
		errors.Panic("用户名不合法，请不要使用非法字符")
	}

	var u *model.User
	var userId int64
	if g.Config.Ldap.Enabled {
		sucess, err := utils.LdapBind(g.Config.Ldap.Addr,
			g.Config.Ldap.BaseDN,
			g.Config.Ldap.BindDN,
			g.Config.Ldap.BindPasswd,
			g.Config.Ldap.UserField,
			username,
			password)

		errors.MaybePanic(err)
		if !sucess {
			errors.Panic("name or password error")
			return
		}

		userAttributes, err := utils.Ldapsearch(g.Config.Ldap.Addr,
			g.Config.Ldap.BaseDN,
			g.Config.Ldap.BindDN,
			g.Config.Ldap.BindPasswd,
			g.Config.Ldap.UserField,
			username,
			g.Config.Ldap.Attributes)
		userSn := ""
		userMail := ""
		userTel := ""
		if err == nil {
			userSn = userAttributes["sn"]
			userMail = userAttributes["mail"]
			userTel = userAttributes["telephoneNumber"]
		}

		arr := strings.Split(username, "@")
		var userName, userEmail string
		if len(arr) == 2 {
			userName = arr[0]
			userEmail = username
		} else {
			userName = username
			userEmail = userMail
		}

		u, err = model.GetUserByName(userName)
		errors.MaybePanic(err)
		if u == nil {
			// 说明用户不存在
			u = &model.User{
				Name:     userName,
				Password: "",
				Cnname:   userSn,
				Phone:    userTel,
				Email:    userEmail,
			}
			errors.MaybePanic(u.Save())
		}
		userId = u.Id
	} else {
		var err error
		userId, err = model.UserLogin(username, utils.EncryptPassword(g.Config.Salt, password))
		errors.MaybePanic(err)
	}

	render.Data(w, cookie.WriteUser(w, userId, username))
}

func ApiLogin(w http.ResponseWriter, r *http.Request) {
	_, name, _ :=cookie.ReadUser(r)
	result, _ := ioutil.ReadAll(r.Body)
	defer  r.Body.Close()
	var user map[string]string
	err := json.Unmarshal(result, &user)
	if err != nil {
		errors.MaybePanic(err)
		render.ErrorCode(w, err)
		return
	}
	username := user["username"]
	password := user["password"]

	if name != "" && username == name {
		render.ErrorCode(w, errors.NewError(username + "已经登录"))
		return
	}
	if str.HasDangerousCharacters(username) {
		errors.Panic("用户名不合法，请不要使用非法字符")
		render.ErrorCode(w, errors.NewError("用户名不合法，请不要使用非法字符"))
		return
	}

	userId, err := model.UserLogin(username, utils.EncryptPassword(g.Config.Salt, password))
	if err != nil {
		render.ErrorCode(w, err)
		return
	}
	// 写入cookie
	if err != cookie.WriteUser(w, userId, username) {
		render.ErrorCode(w, err)
		return
	}
	render.Data(w, "", "登录成功")
}

func MeJson(w http.ResponseWriter, r *http.Request) {
	render.Data(w, MeRequired(LoginRequired(w, r)))
}

func UsersJson(w http.ResponseWriter, r *http.Request) {
	MeRequired(LoginRequired(w, r))
	query := param.String(r, "query", "")
	limit := param.Int(r, "limit", 10)
	if str.HasDangerousCharacters(query) {
		errors.Panic("query invalid")
		return
	}

	users, err := model.QueryUsers(query, limit)
	errors.MaybePanic(err)

	render.Data(w, users)
}

func UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	me := MeRequired(LoginRequired(w, r))

	cnname := param.String(r, "cnname", "")
	email := param.String(r, "email", "")
	phone := param.String(r, "phone", "")
	wechat := param.String(r, "wechat", "")

	if str.HasDangerousCharacters(cnname) {
		errors.Panic("中文名不合法")
	}
	if email != "" && !str.IsMail(email) {
		errors.Panic("邮箱不合法")
	}
	if phone != "" && !str.IsPhone(phone) {
		errors.Panic("手机号不合法")
	}
	if str.HasDangerousCharacters(wechat) {
		errors.Panic("微信不合法")
	}

	me.Cnname = cnname
	me.Email = email
	me.Phone = phone
	me.Wechat = wechat
	errors.MaybePanic(me.UpdateProfile())
	render.Data(w, "ok")
}

func ChangeMyPasswd(w http.ResponseWriter, r *http.Request) {

	uid, _ := LoginRequired(w, r)
	me, err := model.GetUserPwById(uid)
	errors.MaybePanic(err)

	oldPasswd := param.MustString(r, "old_password")
	newPasswd := param.MustString(r, "new_password")
	repeat := param.MustString(r, "repeat")

	if newPasswd != repeat {
		errors.Panic("两次输入的密码不一致")
	}

	err = me.ChangePasswd(utils.EncryptPassword(g.Config.Salt, oldPasswd), utils.EncryptPassword(g.Config.Salt, newPasswd))
	if err == nil {
		cookie.RemoveUser(w)
	}

	errors.MaybePanic(err)
	render.Data(w, "ok")
}
