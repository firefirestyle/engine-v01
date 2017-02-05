package handler

import (
	"net/http"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	"github.com/firefirestyle/engine-v01/user/user"
	"google.golang.org/appengine"
)

func (obj *UserHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	cursor := values.Get("cursor")
	mode := values.Get("mode")
	keyOnly := true
	var foundObj *user.FoundUser = nil
	if mode == "-point" {
		foundObj = obj.manager.FindUserWithPoint(ctx, cursor, keyOnly)
	} else {
		foundObj = obj.manager.FindUserWithNewOrder(ctx, cursor, keyOnly)
	}
	propObj.SetPropStringList("", "keys", foundObj.UserIds)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorNext", foundObj.CursorNext)

	w.Write(propObj.ToJson())
}
