package handler

import (
	"net/http"

	"io/ioutil"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleUpdateInfo(w http.ResponseWriter, r *http.Request) {
	outputProp := miniprop.NewMiniProp()
	v, _ := ioutil.ReadAll(r.Body)
	inputProp := miniprop.NewMiniPropFromJson(v)
	ctx := appengine.NewContext(r)
	userName := inputProp.GetString("userName", "")
	displayName := inputProp.GetString("displayName", "")
	content := inputProp.GetString("content", "")
	token := inputProp.GetString("token", "")

	//
	// check token
	{
		loginResult := obj.CheckLoginFromToken(r, token, false)
		if loginResult.IsLogin == false {
			obj.HandleError(w, r, miniprop.NewMiniProp(), 2001, "need to login")
			return
		}

		if userName == "" {
			userName = loginResult.AccessTokenObj.GetUserName()
		}
		if userName != loginResult.AccessTokenObj.GetUserName() {
			usrObj, userErr := obj.GetManager().GetUserFromUserName(ctx, loginResult.AccessTokenObj.GetUserName())
			if userErr != nil {
				obj.HandleError(w, r, outputProp, 2002, userErr.Error())
				return
			}
			if true == usrObj.IsMaster() {
				obj.HandleError(w, r, outputProp, 2002, "need to admin status ")
			}
		}
	}

	usrObj, userErr := obj.GetManager().GetUserFromUserName(ctx, userName)
	if userErr != nil {
		obj.HandleError(w, r, outputProp, 2002, userErr.Error())
		return
	}
	usrObj.SetDisplayName(displayName)
	usrObj.SetCont(content)
	nextUserObj, nextUserErr := obj.GetManager().SaveUserWithImmutable(ctx, usrObj)
	if nextUserErr != nil {
		obj.HandleError(w, r, outputProp, 2004, userErr.Error())
		return
	}
	//
	outputProp.CopiedOver(miniprop.NewMiniPropFromMap(nextUserObj.ToMapPublic()))
	w.WriteHeader(http.StatusOK)
	w.Write(outputProp.ToJson())

}
