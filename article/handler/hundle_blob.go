package handler

import (
	"net/http"

	"strings"

	"io/ioutil"

	miniprop "github.com/firefirestyle/engine-v01/prop"
)

func (obj *ArticleHandler) GetArticleIdFromDir(dir string) string {
	if false == strings.HasPrefix(dir, "/art/") {
		return ""
	}
	t1 := strings.Replace(dir, "/art/", "", 1)
	t2 := strings.Index(t1, "/")
	if t2 == -1 {
		t2 = len(t1)
	}

	return t1[0:t2]
}

func (obj *ArticleHandler) GetDirFromDir(dir string) string {
	if false == strings.HasPrefix(dir, "/art/") {
		return ""
	}
	t1 := strings.Replace(dir, "/art/", "", 1)
	t2 := strings.Index(t1, "/")
	if t2 == -1 {
		return ""
	}

	return t1[t2:len(t1)]
}

func (obj *ArticleHandler) MakePath(articleId, dir string) string {
	if strings.HasPrefix(dir, "/") {
		dir = dir[1:]
	}
	t1 := "/art/" + articleId + "/" + dir
	if strings.HasSuffix(t1, "/") {
		t1 = t1[0:(len(t1) - 1)]
	}
	return t1
}

func (obj *ArticleHandler) HandleBlobRequestToken(w http.ResponseWriter, r *http.Request) {
	//
	// load param from json
	params, _ := ioutil.ReadAll(r.Body)
	inputPropObj := miniprop.NewMiniPropFromJson(params)
	dir := inputPropObj.GetString("dir", "")
	name := inputPropObj.GetString("file", "")
	articleId := inputPropObj.GetString("articleId", "")
	t1 := obj.MakePath(articleId, dir)
	obj.blobHundler.HandleBlobRequestTokenFromParams(w, r, t1, name, inputPropObj)
}

func (obj *ArticleHandler) HandleBlobUpdated(w http.ResponseWriter, r *http.Request) {
	obj.blobHundler.HandleUploaded(w, r)
}

func (obj *ArticleHandler) HandleBlobGet(w http.ResponseWriter, r *http.Request) {
	obj.blobHundler.HandleGet(w, r)
}
