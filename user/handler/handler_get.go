package handler

import (
	"net/http"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	minisession "github.com/firefirestyle/engine-v01/session"
	miniuser "github.com/firefirestyle/engine-v01/user/user"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	userName := values.Get("userName")
	sign := values.Get("sign")
	key := values.Get("key")
	obj.HandleGetBase(w, r, userName, sign, key, false)
}

func (obj *UserHandler) HandleGetMe(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	inputProp := miniprop.NewMiniPropFromJsonReader(r.Body)
	token := inputProp.GetString("token", "")
	loginResult := obj.GetSessionMgr().CheckAccessToken(ctx, token, minisession.MakeOptionInfo(r), true)
	userName := loginResult.AccessTokenObj.GetUserName()
	if loginResult.IsLogin == false {
		userName = ""
	}
	obj.HandleGetBase(w, r, userName, "", "", false)
}

func (obj *UserHandler) HandleGetBase(w http.ResponseWriter, r *http.Request, userName string, sign string, key string, includePrivate bool) {
	ctx := appengine.NewContext(r)
	var usrObj *miniuser.User = nil
	var userErr error = nil

	outputProp := miniprop.NewMiniProp()
	if userName != "" {
		usrObj, userErr = obj.GetManager().GetUserFromUserName(ctx, userName)
	} else if key != "" {
		usrObj, userErr = obj.GetManager().GetUserFromKey(ctx, key)
	} else {
		obj.HandleError(w, r, outputProp, 2002, "wrong request")
		return
	}

	if userErr != nil {
		obj.HandleError(w, r, outputProp, 2002, userErr.Error())
		return
	}
	//
	//
	if key != "" || sign != "" {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
	}

	if includePrivate == true {
		outputProp.CopiedOver(miniprop.NewMiniPropFromMap(usrObj.ToMapAll()))
	} else {
		outputProp.CopiedOver(miniprop.NewMiniPropFromMap(usrObj.ToMapPublic()))
	}
	Debug(ctx, "--cont-- "+usrObj.GetCont())
	w.Write(outputProp.ToJson())
	return
}

func (obj *UserHandler) CheckLogin(r *http.Request, token string, useIp bool) minisession.CheckResult {
	return obj.GetSessionMgr().CheckAccessToken(appengine.NewContext(r), token, minisession.MakeOptionInfo(r), useIp)
}
