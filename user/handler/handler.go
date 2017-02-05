package handler

import (
	"net/http"

	"crypto/sha1"

	miniblob "github.com/firefirestyle/engine-v01/blob/blob"
	blobhandler "github.com/firefirestyle/engine-v01/blob/handler"
	"github.com/firefirestyle/engine-v01/oauth/twitter"
	miniprop "github.com/firefirestyle/engine-v01/prop"
	minisession "github.com/firefirestyle/engine-v01/session"
	miniuser "github.com/firefirestyle/engine-v01/user/user"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

type UserHandler struct {
	manager        *miniuser.UserManager
	sessionMgr     *minisession.SessionManager
	blobHandler    *blobhandler.BlobHandler
	twitterHandler *twitter.TwitterHandler
	onEvents       UserHandlerOnEvent
	completeFunc   func(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, hh *blobhandler.BlobHandler, blobObj *miniblob.BlobItem) error
}

type UserHandlerManagerConfig struct {
	UserKind                   string
	SessionKind                string
	BlobKind                   string
	BlobPointerKind            string
	BlobSign                   string
	MemcachedOnlyInBlobPointer bool
	LengthHash                 int
}

type UserHandlerOnEvent struct {
	OnGetUserRequestList       []func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp) error
	OnGetUserFailedList        []func(w http.ResponseWriter, r *http.Request, h *UserHandler, o *miniprop.MiniProp)
	OnGetUserSuccessList       []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniuser.User, o *miniprop.MiniProp) error
	OnUpdateUserRequestList    []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error
	OnUpdateUserFailedList     []func(w http.ResponseWriter, r *http.Request, h *UserHandler, i *miniprop.MiniProp, o *miniprop.MiniProp)
	OnUpdateUserBeforeSaveList []func(w http.ResponseWriter, r *http.Request, h *UserHandler, u *miniuser.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error
	OnUpdateUserSuccessList    []func(w http.ResponseWriter, r *http.Request, h *UserHandler, u *miniuser.User, i *miniprop.MiniProp, o *miniprop.MiniProp) error
}

func NewUserHandler(callbackUrl string, //
	config UserHandlerManagerConfig) *UserHandler {
	if config.UserKind == "" {
		config.UserKind = "fu"
	}
	if config.SessionKind == "" {
		config.SessionKind = config.UserKind + "-session"
	}
	if config.BlobKind == "" {
		config.BlobKind = config.UserKind + "-blob"
	}
	if config.BlobPointerKind == "" {
		config.BlobPointerKind = config.UserKind + "-blob-pointer"
	}
	if config.BlobSign == "" {
		config.BlobSign = string(sha1.New().Sum([]byte("" + config.UserKind)))
	}
	//

	ret := &UserHandler{
		manager: miniuser.NewUserManager(miniuser.UserManagerConfig{
			UserKind:   config.UserKind,
			LengthHash: config.LengthHash,
		}),
		sessionMgr: minisession.NewSessionManager(minisession.SessionManagerConfig{
			Kind: config.SessionKind,
		}),
		blobHandler: blobhandler.NewBlobHandler(callbackUrl, config.BlobSign, miniblob.BlobManagerConfig{
			Kind:        config.BlobKind,
			PointerKind: config.BlobPointerKind,
			CallbackUrl: callbackUrl,
			HashLength:  10,
		}),
		onEvents: UserHandlerOnEvent{},
	}

	ret.blobHandler.AddOnBlobComplete(ret.OnBlobComplete)
	return ret
}

func (obj *UserHandler) GetBlobHandler() *blobhandler.BlobHandler {
	return obj.blobHandler
}

func (obj *UserHandler) AddTwitterSession(twitterConfig twitter.TwitterOAuthConfig) {
	obj.twitterHandler = obj.NewTwitterHandlerObj(twitterConfig)
}

func (obj *UserHandler) GetSessionMgr() *minisession.SessionManager {
	return obj.sessionMgr
}

func (obj *UserHandler) GetManager() *miniuser.UserManager {
	return obj.manager
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

func (obj *UserHandler) CheckLoginFromToken(r *http.Request, token string, useIp bool) minisession.CheckResult {
	ctx := appengine.NewContext(r)
	return obj.GetSessionMgr().CheckAccessToken(ctx, token, minisession.MakeOptionInfo(r), useIp)
}

func (obj *UserHandler) HandleError(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, errorCode int, errorMessage string) {
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
