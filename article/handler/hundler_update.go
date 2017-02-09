package handler

import (
	"net/http"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	obj.HandleUpdateBase(w, r, obj.GetInputProp(w, r))
}

func (obj *ArticleHandler) HandleUpdateBase(w http.ResponseWriter, r *http.Request, inputProp *miniprop.MiniProp) {
	ctx := appengine.NewContext(r)
	propObj := miniprop.NewMiniProp()
	//
	// load param from json

	articleId := inputProp.GetString("articleId", "")
	ownerName := inputProp.GetString("userName", "")

	//
	//
	//
	if articleId == "" {
		obj.HandleError(w, r, miniprop.NewMiniProp(), ErrorCodeNotFoundArticleId, "Not Found Article")
		return
	}

	artObj, errGetArt := obj.GetManager().GetArticleFromPointer(ctx, articleId)
	if errGetArt != nil {
		obj.HandleError(w, r, miniprop.NewMiniProp(), ErrorCodeNotFoundArticleId, "Not Found Article")
		return
	}
	//
	if inputProp.Contain("title") {
		title := inputProp.GetString("title", "")
		artObj.SetTitle(title)
	}

	if inputProp.Contain("userName") {
		artObj.SetUserName(ownerName)
	}

	if inputProp.Contain("content") {
		content := inputProp.GetString("content", "")
		artObj.SetCont(content)
	}

	if inputProp.Contain("info") {
		content := inputProp.GetString("info", "")
		artObj.SetInfo(content)
	}

	if inputProp.Contain("tags") {
		tags := inputProp.GetPropStringList("", "tags", make([]string, 0))
		artObj.SetTags(tags)
	}

	if inputProp.Contain("lat") {
		lat := inputProp.GetFloat("lat", -999.0)
		artObj.SetLat(lat)
	}

	if inputProp.Contain("lng") {
		lng := inputProp.GetFloat("lng", -999.0)
		artObj.SetLng(lng)
	}

	if inputProp.Contain("iconUrl") {
		iconUrl := inputProp.GetString("iconUrl", "")
		artObj.SetIconUrl(iconUrl)
	}

	//
	//
	if inputProp.Contain("propKeys") {
		propKeys := inputProp.GetPropStringList("", "propKeys", make([]string, 0))
		propValues := inputProp.GetPropStringList("", "propValues", make([]string, 0))
		artObj.ClearProp()
		if len(propKeys) == len(propValues) {
			for i, kv := range propKeys {
				artObj.SetProp(kv, propValues[i])
			}
		}
	}
	//
	//
	_, errSave := obj.GetManager().SaveArticleWithImmutable(ctx, artObj)

	if errSave != nil {
		obj.HandleError(w, r, miniprop.NewMiniProp(), ErrorCodeFailedToSave, errSave.Error())
		return
	} else {
		propObj.SetPropString("", "articleId", artObj.GetArticleId())
		propObj.SetPropString("", "articleKey", artObj.GetStringId())
		w.WriteHeader(http.StatusOK)
		w.Write(propObj.ToJson())
	}
}
