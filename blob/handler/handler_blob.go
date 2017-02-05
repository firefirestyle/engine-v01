package handler

import (
	//	"net/url"

	"io/ioutil"
	"net/http"

	mm "github.com/firefirestyle/engine-v01/prop"

	"strconv"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
)

func (obj *BlobHandler) HandleBlobRequestToken(w http.ResponseWriter, r *http.Request) {
	params, _ := ioutil.ReadAll(r.Body)
	inputPropObj := mm.NewMiniPropFromJson(params)
	dirName := inputPropObj.GetString("dir", "")
	fileName := inputPropObj.GetString("file", "")
	obj.HandleBlobRequestTokenFromParams(w, r, dirName, fileName, inputPropObj)
}

func (obj *BlobHandler) HandleBlobRequestTokenFromParams(w http.ResponseWriter, r *http.Request, dirName string, fileName string, inputPropObj *mm.MiniProp) {
	ctx := appengine.NewContext(r)
	outputPropObj := mm.NewMiniProp()
	if inputPropObj == nil {
		params, _ := ioutil.ReadAll(r.Body)
		inputPropObj = mm.NewMiniPropFromJson(params)
	}

	//
	//
	vs := map[string]string{}
	//
	//
	kv := strconv.FormatInt(time.Now().Unix(), 36)
	reqUrl, reqName, err := obj.manager.MakeRequestUrl(ctx, dirName, fileName, kv, obj.privateSign, vs)
	if err != nil {
		HandleError(w, r, outputPropObj, ErrorCodeAtBlobMakeRequestUrl, "failed to make uploadurl")
	} else {
		outputPropObj.SetString("token", reqUrl.String())
		outputPropObj.SetString("name", reqName)
		w.Write(outputPropObj.ToJson())
		w.WriteHeader(http.StatusOK)
	}
}

func (obj *BlobHandler) HandleUploaded(w http.ResponseWriter, r *http.Request) {
	//
	//
	outputPropObj := mm.NewMiniProp()
	res, e := obj.manager.CheckedCallback(r, obj.privateSign)
	if e != nil {
		HandleError(w, r, outputPropObj, ErrorCodeAtBlobCheckCallback, e.Error())
		return
	}
	curTime := time.Now().Unix()
	kvTime, errTime := strconv.ParseInt(r.FormValue("kv"), 36, 64)
	if errTime != nil || !(curTime-60*1 < kvTime && kvTime < curTime+60*10) {
		HandleError(w, r, outputPropObj, ErrorCodeAtBlobCheckCallback, "kv time error")
		return
	}
	Debug(appengine.NewContext(r), "(3) ---- ")
	//
	ctx := appengine.NewContext(r)
	newItem := obj.manager.NewBlobItem(ctx, res.DirName, res.FileName, res.BlobKey)
	//
	err2 := obj.manager.SaveBlobItemWithImmutable(ctx, newItem)
	if err2 != nil {
		blobstore.Delete(ctx, appengine.BlobKey(res.BlobKey))
		HandleError(w, r, outputPropObj, ErrorCodeAtBlobSaveBlobItem, err2.Error())
		return
	}
	Debug(appengine.NewContext(r), "(4) ---- ")
	err3 := obj.OnBlobComplete(w, r, outputPropObj, obj, newItem)
	if err3 != nil {
		obj.GetManager().DeleteBlobItem(ctx, newItem)
		HandleError(w, r, outputPropObj, ErrorCodeAtBlobCompleteCheck, err3.Error())
		return
	}
	outputPropObj.SetString("blobkey", newItem.GetBlobKey())
	w.Write(outputPropObj.ToJson())
	w.WriteHeader(http.StatusOK)
}
