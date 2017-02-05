package twitter

//
// https://dev.twitter.com/oauth/overview
// https://dev.twitter.com/web/sign-in/implementing
//
import (
	"errors"
	"net/url"

	"golang.org/x/net/context"
)

const (
	RequestTokenURl        = "https://api.twitter.com/oauth/request_token"
	AccessTokenURL         = "https://api.twitter.com/oauth/access_token"
	OAuthToken             = "oauth_token"
	OAuthTokenSecret       = "oauth_token_secret"
	OAuthCallbackConfirmed = "oauth_callback_confirmed"
	OAuthVerifier          = "oauth_verifier"
	UserID                 = "user_id"
	ScreenName             = "screen_name"
)

type TwitterManager struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	AllowInvalidSSL   bool
}

type Twitter struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	AllowInvalidSSL   bool
	//oauthObj          *OAuth1Client
}

func NewTwitterManager(consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string, allowInvalidSSL bool) *TwitterManager {
	ret := new(TwitterManager)
	ret.ConsumerKey = consumerKey
	ret.ConsumerSecret = consumerSecret
	ret.AccessToken = accessToken
	ret.AccessTokenSecret = accessTokenSecret
	ret.AllowInvalidSSL = allowInvalidSSL
	return ret
}

func (obj *TwitterManager) NewTwitter() *Twitter {
	ret := new(Twitter)
	ret.ConsumerKey = obj.ConsumerKey
	ret.ConsumerSecret = obj.ConsumerSecret
	ret.AccessToken = obj.AccessToken
	ret.AccessTokenSecret = obj.AccessTokenSecret
	ret.AllowInvalidSSL = obj.AllowInvalidSSL
	return ret
}

//
// OAuthToken
// OAuthTokenSecret
// OAuthCallbackConfirmed
func (obj *Twitter) SendRequestToken(ctx context.Context, callbackUrl string) (*SendRequestTokenResult, error) {
	//
	//
	//
	oauthObj := NewOAuthClient(obj.ConsumerKey, obj.ConsumerSecret, obj.AccessToken, obj.AccessTokenSecret, obj.AllowInvalidSSL)
	oauthObj.Callback = callbackUrl
	result, err := oauthObj.Post(ctx, RequestTokenURl, make(map[string]string, 0), "")
	if err != nil {
		return nil, err
	}
	//
	//
	//
	keyValueObj := NewSendRequestTokenResult(result)
	oauth_token := keyValueObj.GetOAuthToken()
	if oauth_token == "" {
		return nil, errors.New("oauthtoken is nil")
	}

	return keyValueObj, nil
}

//
// OAuthToken
// OAuthTokenSecret
// UserID
// ScreenName
func (obj *Twitter) OnCallbackSendRequestToken(ctx context.Context, url *url.URL) (*SendAccessTokenResult, error) {
	q := url.Query()
	verifiers := q[OAuthVerifier]
	tokens := q[OAuthToken]

	if len(verifiers) != 1 || len(tokens) != 1 {
		return nil, errors.New("unexpected query")
	}
	ret1 := make(map[string]string, 0)
	ret1[OAuthVerifier] = verifiers[0]
	ret1[OAuthToken] = tokens[0]
	ret2, ret3 := obj.SendAccessToken(ctx, tokens[0], verifiers[0])
	return ret2, ret3
}

//
// OAuthToken
// OAuthTokenSecret
// UserID
// ScreenName
func (obj *Twitter) SendAccessToken(ctx context.Context, oauthToken string, oauthVerifier string) (*SendAccessTokenResult, error) {
	//
	//
	//
	oauthObj := NewOAuthClient(obj.ConsumerKey, obj.ConsumerSecret, obj.AccessToken, obj.AccessTokenSecret, obj.AllowInvalidSSL)
	oauthObj.Callback = ""
	oauthObj.AccessToken = oauthToken
	result, err := oauthObj.Post( //
		ctx,            //
		AccessTokenURL, //
		map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		"oauth_verifier="+oauthVerifier+"\r\n")
	if err != nil {
		return nil, err
	}
	//
	//
	//
	keyValueObj := NewSendAccessTokenResult(result)
	return keyValueObj, nil
}
