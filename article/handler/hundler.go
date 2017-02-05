package handler

import (
	"net/http"

	"io/ioutil"

	"github.com/firefirestyle/engine-v01/article/article"
	miniblob "github.com/firefirestyle/engine-v01/blob/blob"
	blobhandler "github.com/firefirestyle/engine-v01/blob/handler"
	miniprop "github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const (
	ErrorCodeFailedToSave                = 2001
	ErrorCodeFailedToCheckAboutGetCalled = 2002
	ErrorCodeNotFoundArticleId           = 2003
)

type ArticleHandler struct {
	articleKind string
	blobKind    string
	pointerKind string
	artMana     *article.ArticleManager
	blobHundler *blobhandler.BlobHandler
}

type ArticleHandlerConfig struct {
	ArticleKind     string
	PointerKind     string
	BlobKind        string
	BlobPointerKind string
	TagKind         string
	BlobCallbackUrl string
	BlobSign        string
	MemcachedOnly   bool
	LengthHash      int
}

func NewArtHandler(config ArticleHandlerConfig) *ArticleHandler {
	if config.ArticleKind == "" {
		config.ArticleKind = "ffart"
	}
	if config.PointerKind == "" {
		config.PointerKind = config.ArticleKind + "-pointer"
	}
	if config.BlobKind == "" {
		config.BlobKind = config.ArticleKind + "-blob"
	}
	if config.BlobPointerKind == "" {
		config.BlobPointerKind = config.ArticleKind + "-blob-pointer"
	}
	if config.TagKind == "" {
		config.TagKind = config.ArticleKind + "-tag"
	}
	artMana := article.NewArticleManager(article.ArticleManagerConfig{
		KindArticle:    config.ArticleKind,
		KindPointer:    config.PointerKind,
		LimitOfFinding: 20,
		LengthHash:     config.LengthHash,
	})
	//
	//
	artHandlerObj := &ArticleHandler{
		articleKind: config.ArticleKind,
		blobKind:    config.BlobKind,
		artMana:     artMana,
	}

	//
	artHandlerObj.blobHundler = blobhandler.NewBlobHandler(config.BlobCallbackUrl, config.BlobSign,
		miniblob.BlobManagerConfig{
			Kind:                   config.BlobKind,
			CallbackUrl:            config.BlobCallbackUrl,
			PointerKind:            config.BlobPointerKind,
			MemcachedOnlyInPointer: config.MemcachedOnly,
			HashLength:             10,
		})
	artHandlerObj.blobHundler.AddOnBlobComplete(func(w http.ResponseWriter, r *http.Request, o *miniprop.MiniProp, hh *blobhandler.BlobHandler, i *miniblob.BlobItem) error {
		dirSrc := r.URL.Query().Get("dir")
		articlId := artHandlerObj.GetArticleIdFromDir(dirSrc)
		dir := artHandlerObj.GetDirFromDir(dirSrc)
		fileName := r.URL.Query().Get("file")
		//
		//
		if dir == "" && fileName == "icon" {
			ctx := appengine.NewContext(r)
			artObj, errGet := artHandlerObj.GetManager().GetArticleFromPointer(ctx, articlId)
			if errGet != nil {
				return errGet
			}

			artObj.SetIconUrl("key://" + i.GetBlobKey())
			_, errSave := artHandlerObj.GetManager().SaveArticleWithImmutable(ctx, artObj)

			if errSave != nil {
				return errSave
			}
		}
		return nil
	})
	return artHandlerObj
}

func (obj *ArticleHandler) GetManager() *article.ArticleManager {
	return obj.artMana
}

func (obj *ArticleHandler) GetBlobHandler() *blobhandler.BlobHandler {
	return obj.blobHundler
}

func (obj *ArticleHandler) HandleError(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, errorCode int, errorMessage string) {
	//
	//
	if outputProp == nil {
		outputProp = miniprop.NewMiniProp()
	}
	if errorCode != 0 {
		outputProp.SetInt("errorCode", errorCode)
	}
	if errorMessage != "" {
		outputProp.SetString("errorMessage", errorMessage)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(outputProp.ToJson())
}

func (obj *ArticleHandler) GetInputProp(w http.ResponseWriter, r *http.Request) *miniprop.MiniProp {
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return miniprop.NewMiniProp()
	} else {
		return miniprop.NewMiniPropFromJson(v)
	}
}

//
//
//

// HandleBlobRequestTokenFromParams

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

///
//
