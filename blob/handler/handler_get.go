package handler

import (
	//	"net/url"

	"net/http"

	mm "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
)

func (obj *BlobHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	requestValues := r.URL.Query()
	key := requestValues.Get("key")
	dir := requestValues.Get("dir")
	file := requestValues.Get("file")

	obj.HandleGetBase(w, r, key, dir, file)
}

func (obj *BlobHandler) HandleGetBase(w http.ResponseWriter, r *http.Request, key, dir, file string) {

	//
	outputPropObj := mm.NewMiniProp()
	//
	if key != "" {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
		blobstore.Send(w, appengine.BlobKey(key))
		return
	}
	//
	ctx := appengine.NewContext(r)
	blobObj, err := obj.manager.GetBlobItemFromPointer(ctx, dir, file)
	if err != nil {
		HandleError(w, r, outputPropObj, ErrorCodeAtGetRequestFindBlobItem, err.Error())
		return
	} else {
		blobstore.Send(w, appengine.BlobKey(blobObj.GetBlobKey()))
		return
	}

}
