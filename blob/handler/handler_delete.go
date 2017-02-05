package handler

import (
	"net/http"

	mm "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
)

func (obj *BlobHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	outputPropObj := mm.NewMiniProp()
	requestValues := r.URL.Query()
	dir := requestValues.Get("dir")
	file := requestValues.Get("file")
	ctx := appengine.NewContext(r)
	//
	//

	//
	//
	blolStringId, _, err := obj.manager.GetBlobItemStringIdFromPointer(ctx, dir, file)
	//	obj.manager.DeleteBlobItemWithPointerFromStringId(ctx, obj.manager.MakeStringId())
	//	blobObj, _, err := obj.manager.GetBlobItemFromPointer(ctx, dir, file)
	if err != nil {
		HandleError(w, r, outputPropObj, ErrorCodeAtDeleteRequestFindBlobItem, err.Error())
		return
	}
	errDelete := obj.manager.DeleteBlobItemWithPointerFromStringId(ctx, blolStringId)
	if errDelete != nil {
		HandleError(w, r, outputPropObj, ErrorCodeAtDeleteRequestDeleteBlobItem, err.Error())
		return
	} else {
		w.Write(outputPropObj.ToJson())
		return
	}
}
