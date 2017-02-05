package handler

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	miniblob "github.com/firefirestyle/engine-v01/blob/blob"
	blobhandler "github.com/firefirestyle/engine-v01/blob/handler"
	miniprop "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
)

func (obj *UserHandler) GetUserNameFromDir(dir string) string {
	if false == strings.HasPrefix(dir, "/user/") {
		return ""
	}
	t1 := strings.Replace(dir, "/user/", "", 1)
	t2 := strings.Index(t1, "/")
	if t2 == -1 {
		t2 = len(t1)
	}

	return t1[0:t2]
}

func (obj *UserHandler) GetDirFromDir(dir string) string {
	if false == strings.HasPrefix(dir, "/user/") {
		return ""
	}
	t1 := strings.Replace(dir, "/user/", "", 1)
	t2 := strings.Index(t1, "/")
	if t2 == -1 {
		t2 = 0
	}

	return t1[t2:len(t1)]
}

func (obj *UserHandler) MakeDir(userName string, dir string) string {
	if dir == "" {
		return "/user/" + userName
	}
	if strings.HasPrefix(dir, "/") {
		dir = dir[1:]
	}
	if strings.HasSuffix(dir, "/") {
		dir = dir[0 : len(dir)-1]
	}

	return "/user/" + userName + "/" + dir
}

func (obj *UserHandler) HandleBlobRequestToken(w http.ResponseWriter, r *http.Request) {
	params, _ := ioutil.ReadAll(r.Body)
	obj.HandleBlobRequestTokenBase(w,r,miniprop.NewMiniPropFromJson(params))
}

func (obj *UserHandler) HandleBlobRequestTokenBase(w http.ResponseWriter, r *http.Request,inputPropObj *miniprop.MiniProp) {
	token := inputPropObj.GetString("token", "")
	if token == "" {
		obj.HandleError(w, r, miniprop.NewMiniProp(), 2001, "not found token")
		return
	}
	loginResult := obj.CheckLoginFromToken(r, token, false)
	if loginResult.IsLogin == false {
		obj.HandleError(w, r, miniprop.NewMiniProp(), 2001, "need to login")
		return
	}
	userName := loginResult.AccessTokenObj.GetUserName()
	dir := inputPropObj.GetString("dir", "")
	name := inputPropObj.GetString("file", "")

	//
	//
	obj.blobHandler.HandleBlobRequestTokenFromParams(w, r, obj.MakeDir(userName, dir), name, inputPropObj)
}

func (obj *UserHandler) HandleBlobUpdated(w http.ResponseWriter, r *http.Request) {
	//
	obj.blobHandler.HandleUploaded(w, r)
}

func (obj *UserHandler) HandleBlobGet(w http.ResponseWriter, r *http.Request) {
	//
	values := r.URL.Query()
	key := values.Get("key")
	dir := values.Get("dir")
	file := values.Get("file")
	userName := values.Get("userName")

	obj.blobHandler.HandleGetBase(w, r, key, obj.MakeDir(userName, dir), file)
}

func (userMgrObj *UserHandler) OnBlobComplete(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, hh *blobhandler.BlobHandler, blobObj *miniblob.BlobItem) error {
	dir := r.URL.Query().Get("dir")
	if true == strings.HasPrefix(dir, "/user") {
		ctx := appengine.NewContext(r)
		userName := userMgrObj.GetUserNameFromDir(dir)

		userObj, userErr := userMgrObj.GetManager().GetUserFromUserName(ctx, userName)
		if userErr != nil {
			outputProp.SetString("error", "not found user")
			return userErr
		}
		userObj.SetIconUrl("key://" + blobObj.GetBlobKey())
		userMgrObj.GetManager().SaveUserWithImmutable(ctx, userObj)
		if userMgrObj.completeFunc != nil {
			return userMgrObj.completeFunc(w, r, outputProp, hh, blobObj)
		} else {
			return nil
		}
	} else {
		return errors.New("unsupport")
	}
}
