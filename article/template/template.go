package template

import (
	"net/http"

	arthundler "github.com/firefirestyle/engine-v01/article/handler"
	miniprop "github.com/firefirestyle/engine-v01/prop"
	minisession "github.com/firefirestyle/engine-v01/session"
	userHandler "github.com/firefirestyle/engine-v01/user/handler"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"errors"
)

const (
	UrlArtNew             = "/new"
	UrlArtUpdate          = "/update"
	UrlArtFind            = "/find"
	UrlArtFindMe          = "/find_with_token"
	UrlArtGet             = "/get"
	UrlArtBlobGet         = "/getblob"
	UrlArtRequestBlobUrl  = "/requestbloburl"
	UrlArtCallbackBlobUrl = "/callbackbloburl"
	UrlArtDelete          = "/delete"
)

type ArtTemplateConfig struct {
	KindBaseName string
	PrivateKey   string
	BasePath     string
}

type ArtTemplate struct {
	config         ArtTemplateConfig
	artHandlerObj  *arthundler.ArticleHandler
	getUserHundler func(context.Context) *userHandler.UserHandler
}

func NewArtTemplate(config ArtTemplateConfig, getUserHundler func(context.Context) *userHandler.UserHandler) *ArtTemplate {
	if config.KindBaseName == "" {
		config.KindBaseName = "fa"
	}

	if config.BasePath == "" {
		config.BasePath = "/api/v1/art"
	}

	return &ArtTemplate{
		config:         config,
		getUserHundler: getUserHundler,
	}
}

func (tmpObj *ArtTemplate) CheckArticleOwner(ctx context.Context, checkResult minisession.CheckResult, articleId string) error {
	manager := tmpObj.GetArtHundlerObj(ctx).GetManager()
	artObj, artErr := manager.GetArticleFromPointer(ctx, articleId)
	if artErr != nil {
		return artErr
	}

	if artObj.GetUserName() == checkResult.AccessTokenObj.GetUserName() {
		return nil
	}

	userManager := tmpObj.getUserHundler(ctx).GetManager()
	usrObj, userErr := userManager.GetUserFromUserName(ctx, checkResult.AccessTokenObj.GetUserName())
	if userErr != nil {
		return userErr
	}
	if usrObj.IsMaster() {
		return nil
	}

	return errors.New("failed to check article pwner")
}

func (tmpObj *ArtTemplate) CheckLogin(r *http.Request, token string, useIp bool) minisession.CheckResult {
	check := tmpObj.getUserHundler(appengine.NewContext(r)).CheckLogin(r, token, useIp)
	return check
}

func (tmpObj *ArtTemplate) GetArtHundlerObj(ctx context.Context) *arthundler.ArticleHandler {
	if tmpObj.artHandlerObj == nil {
		tmpObj.artHandlerObj = arthundler.NewArtHandler(
			arthundler.ArticleHandlerConfig{
				ArticleKind:     tmpObj.config.KindBaseName,
				BlobCallbackUrl: tmpObj.config.BasePath + UrlArtCallbackBlobUrl,
				BlobSign:        appengine.VersionID(ctx),
				LengthHash:      10,
			})
	}
	return tmpObj.artHandlerObj
}

func (tmpObj *ArtTemplate) InitArtApi() {

	///
	/// use login check
	///
	http.HandleFunc(tmpObj.config.BasePath+UrlArtNew, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		propObj := miniprop.NewMiniPropFromJsonReader(r.Body)
		loginInfo := tmpObj.CheckLogin(r, propObj.GetString("token", ""), false)
		if loginInfo.IsLogin == false {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4001, "failed to login")
			return
		}
		tmpObj.GetArtHundlerObj(ctx).HandleNewBase(w, r, propObj)

	})

	///
	/// use login check
	///
	http.HandleFunc(tmpObj.config.BasePath+UrlArtUpdate, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		propObj := miniprop.NewMiniPropFromJsonReader(r.Body)
		loginInfo := tmpObj.CheckLogin(r, propObj.GetString("token", ""), false)
		if loginInfo.IsLogin == false {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4001, "failed to login")
			return

		}

		ownerCheckErr := tmpObj.CheckArticleOwner(ctx, loginInfo, propObj.GetString("articleId", ""))
		if ownerCheckErr != nil {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4002, "failed to owner check")
			return
		}

		tmpObj.GetArtHundlerObj(ctx).HandleUpdateBase(w, r, propObj)

	})

	http.HandleFunc(tmpObj.config.BasePath+UrlArtFind, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.GetArtHundlerObj(ctx).HandleFind(w, r)
	})

	///
	/// use login check
	///
	http.HandleFunc(tmpObj.config.BasePath+UrlArtFindMe, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		propObj := miniprop.NewMiniPropFromJsonReader(r.Body)
		loginInfo := tmpObj.CheckLogin(r, propObj.GetString("token", ""), true)
		if loginInfo.IsLogin == false {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4001, "failed to login")
		} else {
			tmpObj.GetArtHundlerObj(ctx).HandleFindBase(w, r, //
				propObj.GetString("cursor", ""), loginInfo.AccessTokenObj.GetUserName(), map[string]string{}, []string{}, "-update")
		}
	})

	http.HandleFunc(tmpObj.config.BasePath+UrlArtGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.GetArtHundlerObj(ctx).HandleGet(w, r)
	})
	//UrlArtGet

	///
	/// use login check
	///
	http.HandleFunc(tmpObj.config.BasePath+UrlArtRequestBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		propObj := miniprop.NewMiniPropFromJsonReader(r.Body)
		loginInfo := tmpObj.CheckLogin(r, propObj.GetString("token", ""), false)
		if loginInfo.IsLogin == false {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4001, "failed to login")
			return
		}
		ownerCheckErr := tmpObj.CheckArticleOwner(ctx, loginInfo, propObj.GetString("articleId", ""))
		if ownerCheckErr != nil {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4002, "failed to owner check")
			return
		}
		tmpObj.GetArtHundlerObj(ctx).HandleBlobRequestTokenBase(w, r, propObj)

	})

	http.HandleFunc(tmpObj.config.BasePath+UrlArtCallbackBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		{
			tmpObj.GetArtHundlerObj(ctx).HandleBlobUpdated(w, r)
		}
	})

	http.HandleFunc(tmpObj.config.BasePath+UrlArtBlobGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.GetArtHundlerObj(ctx).HandleBlobGet(w, r)
	})

	///
	/// use login check
	///
	http.HandleFunc(tmpObj.config.BasePath+UrlArtDelete, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		propObj := miniprop.NewMiniPropFromJsonReader(r.Body)
		//
		//
		loginInfo := tmpObj.CheckLogin(r, propObj.GetString("token", ""), false)
		if loginInfo.IsLogin == false {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4001, "failed to login")
			return
		}
		ownerCheckErr := tmpObj.CheckArticleOwner(ctx, loginInfo, propObj.GetString("articleId", ""))
		if ownerCheckErr != nil {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4002, "failed to owner check")
			return
		}
		tmpObj.GetArtHundlerObj(ctx).HandleDeleteBaseWithFile(w, r, propObj.GetString("articleId", ""), propObj)

	})
	//
}

//
func (tmpObj *ArtTemplate) SaveData(ctx context.Context, token string, articleId string, dirName string, fileName string, data []byte) error {
	return tmpObj.GetArtHundlerObj(ctx).SaveData(ctx, token, tmpObj.config.BasePath+UrlArtRequestBlobUrl, articleId, dirName, fileName, data)
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
