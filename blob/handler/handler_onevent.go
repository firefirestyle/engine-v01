package handler

import (
	"net/http"

	miniblob "github.com/firefirestyle/engine-v01/blob/blob"
	mm "github.com/firefirestyle/engine-v01/prop"
)

func (obj *BlobHandler) AddOnBlobComplete(f func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem) error) {
	obj.onEvent.OnBlobCompleteList = append(obj.onEvent.OnBlobCompleteList, f)
}

func (obj *BlobHandler) OnBlobComplete(w http.ResponseWriter, r *http.Request, o *mm.MiniProp, h *BlobHandler, i *miniblob.BlobItem) error {
	for _, f := range obj.onEvent.OnBlobCompleteList {
		err := f(w, r, o, h, i)
		if err != nil {
			return err
		}
	}
	return nil
}
