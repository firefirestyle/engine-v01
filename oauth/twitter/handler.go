package twitter

import (
	"net/http"
	"net/url"

	"errors"

	"github.com/firefirestyle/engine-v01/oauth/sns"
	"google.golang.org/appengine"

	m "github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const (
	UrlOptCallbackUrl              = "cb"
	UrlOptErrorNotFoundCallbackUrl = "1001"
	UrlOptErrorFailedToMakeToken   = "1002"
)

type TwitterOAuthConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	CallbackUrl       string
	SecretSign        string
	AllowInvalidSSL   bool
}

type TwitterHandler struct {
	twitterManager *TwitterManager
	config         TwitterOAuthConfig
	onEvent        TwitterHundlerOnEvent
}

type TwitterHundlerOnEvent struct {
	OnRequest   func(http.ResponseWriter, *http.Request, *TwitterHandler) (map[string]string, error)
	OnFoundUser func(http.ResponseWriter, *http.Request, *TwitterHandler, *SendAccessTokenResult) map[string]string
}

func NewTwitterHandler( //
	config TwitterOAuthConfig, //
	onEvent TwitterHundlerOnEvent) *TwitterHandler {
	twitterHandlerObj := new(TwitterHandler)
	//	twitterHandlerObj.callbackUrl = callbackUrl
	twitterHandlerObj.twitterManager = NewTwitterManager( //
		config.ConsumerKey, config.ConsumerSecret, config.AccessToken, config.AccessTokenSecret, config.AllowInvalidSSL)
	twitterHandlerObj.config = config

	//
	//
	if onEvent.OnRequest == nil {
		onEvent.OnRequest = func(http.ResponseWriter, *http.Request, *TwitterHandler) (map[string]string, error) {
			return map[string]string{}, nil
		}
	}
	if onEvent.OnFoundUser == nil {
		onEvent.OnFoundUser = func(http.ResponseWriter, *http.Request, *TwitterHandler, *SendAccessTokenResult) map[string]string {
			return map[string]string{}
		}
	}
	twitterHandlerObj.onEvent = onEvent
	return twitterHandlerObj
}

func (obj *TwitterHandler) MakeUrlNotFoundCallbackError(baseAddr string) (string, error) {
	urlObj, err := url.Parse(baseAddr)
	if err != nil {
		return "", err
	}
	query := urlObj.Query()
	query.Add("error", UrlOptErrorNotFoundCallbackUrl)
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}

func (obj *TwitterHandler) MakeUrlFailedToMakeToken(baseAddr string, errMessage string) (string, error) {
	urlObj, err := url.Parse(baseAddr)
	if err != nil {
		return "", err
	}
	query := urlObj.Query()
	query.Add("errorCode", UrlOptErrorFailedToMakeToken)
	query.Add("errorMessage", errMessage)
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}

func (obj *TwitterHandler) HandleLoginEntry(w http.ResponseWriter, r *http.Request) {
	clCallbackUrl := r.URL.Query().Get(UrlOptCallbackUrl)
	// make redirect URL
	if clCallbackUrl == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//
	clCallbackUrlObj, clCallbackUrlErr := url.Parse(clCallbackUrl)
	if clCallbackUrlErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//
	//
	opts, optsErr := obj.onEvent.OnRequest(w, r, obj)
	if optsErr != nil {
		tmpValues := clCallbackUrlObj.Query()
		if opts != nil {
			for k, v := range opts {
				tmpValues.Add(k, v)
			}
		}
		clCallbackUrlObj.RawQuery = tmpValues.Encode()
		http.Redirect(w, r, clCallbackUrlObj.String(), http.StatusFound)
		return
	}
	//
	svCallbackUrlObj, _ := url.Parse(obj.config.CallbackUrl)
	if svCallbackUrlObj.Path == clCallbackUrlObj.Path {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//
	//
	tmpValues := svCallbackUrlObj.Query()
	tmpValues.Add(UrlOptCallbackUrl, clCallbackUrl)

	for k, v := range opts {
		tmpValues.Add(k, v)
	}
	//
	sns.WithHashAndValue(tmpValues, obj.config.SecretSign, clCallbackUrlObj.String(), opts)
	//
	svCallbackUrlObj.RawQuery = tmpValues.Encode()
	//
	//
	redirectUrl := ""

	twitterObj := obj.twitterManager.NewTwitter()
	oauthResult, err := twitterObj.SendRequestToken(appengine.NewContext(r), svCallbackUrlObj.String())
	if err != nil {
		failedOAuthUrl, _ := obj.MakeUrlFailedToMakeToken(clCallbackUrl, err.Error())
		redirectUrl = failedOAuthUrl
	} else {
		redirectUrl = oauthResult.GetOAuthTokenUrl()
	}
	//
	// Do Redirect
	http.Redirect(w, r, redirectUrl, http.StatusFound)

}

func (obj *TwitterHandler) HandleLoginExit(w http.ResponseWriter, r *http.Request) {
	callbackUrl := r.URL.Query().Get(UrlOptCallbackUrl)
	urlObj, err := url.Parse(callbackUrl)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// response easy check
	checkErr := sns.CheckHashAndValue(w, r, obj.config.SecretSign)
	if checkErr != nil {
		query := urlObj.Query()
		query.Add("error", checkErr.Error())
		urlObj.RawQuery = query.Encode()
		http.Redirect(w, r, urlObj.String(), http.StatusFound)
		return
	}

	//
	twitterObj := obj.twitterManager.NewTwitter()
	rt, err := twitterObj.OnCallbackSendRequestToken(appengine.NewContext(r), r.URL)
	if err != nil || rt.GetScreenName() == "" || rt.GetUserID() == "" {
		rt = nil
		if err == nil && (rt.GetScreenName() == "" || rt.GetUserID() == "") {
			err = errors.New("empty user")
		}
	}

	if obj.onEvent.OnFoundUser != nil {
		values := urlObj.Query()
		opts := obj.onEvent.OnFoundUser(w, r, obj, rt)
		for k, v := range opts {
			values.Add(k, v)
		}
		urlObj.RawQuery = values.Encode()
	}
	//

	if err != nil {
		query := urlObj.Query()
		query.Add("error", "oauth")
		urlObj.RawQuery = query.Encode()
		http.Redirect(w, r, urlObj.String(), http.StatusFound)
	} else {
		http.Redirect(w, r, urlObj.String(), http.StatusFound)
	}
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

func HandleError(w http.ResponseWriter, r *http.Request, outputProp *m.MiniProp, errorCode int, errorMessage string) {
	//
	//
	if outputProp == nil {
		outputProp = m.NewMiniProp()
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
