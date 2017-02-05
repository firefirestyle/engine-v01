package handler

import (
	"net/http"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
)

//
// you must to delete file before call this method, if there are many articleid's file.
//
func (obj *ArticleHandler) HandleDeleteBaseWithFile(w http.ResponseWriter, r *http.Request, articleId string, inputObj *miniprop.MiniProp) {
	ctx := appengine.NewContext(r)
	outputObj := miniprop.NewMiniProp()

	deleteFileErr := obj.GetBlobHandler().GetManager().DeleteBlobItemsWithPointerAtRecursiveMode(appengine.NewContext(r), obj.MakePath(articleId, ""))
	if deleteFileErr != nil {
		obj.HandleError(w, r, outputObj, 2002, deleteFileErr.Error())
		return
	}
	err := obj.GetManager().DeleteFromArticleIdWithPointer(ctx, articleId)
	if err != nil {
		obj.HandleError(w, r, outputObj, ErrorCodeNotFoundArticleId, "not found article")
		return
	}
	w.WriteHeader(http.StatusOK)
}
