package handler

import (
	"net/http"

	"github.com/firefirestyle/engine-v01/article/article"
	miniprop "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleGet(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	key := values.Get("key")
	articleId := values.Get("articleId")
	sign := values.Get("sign")
	//	mode := values.Get("m")
	//
	if key != "" {
		keyInfo := obj.GetManager().ExtractInfoFromStringId(key)
		articleId = keyInfo.ArticleId
		sign = keyInfo.Sign
	}
	var artObj *article.Article
	var err error
	//
	//
	if sign != "" {
		artObj, err = obj.GetManager().GetArticleFromArticleId(ctx, articleId, sign)
	} else {
		artObj, err = obj.GetManager().GetArticleFromPointer(ctx, articleId)
	}
	if err != nil {
		obj.HandleError(w, r, miniprop.NewMiniProp(), ErrorCodeNotFoundArticleId, "not found article")
		return
	}
	if sign != "" {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
	}
	propObj := miniprop.NewMiniProp()
	propObj.SetString("articleKey", artObj.GetStringId())
	propObj.CopiedOver(miniprop.NewMiniPropFromMap(artObj.ToMap()))
	w.Write(propObj.ToJson())
}
