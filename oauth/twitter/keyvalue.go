package twitter

import "strings"

type KeyValue struct {
	KeyValues map[string]string
}

func NewKeyValue(content string) *KeyValue {
	ret := new(KeyValue)
	ret.KeyValues = ret.ExtractParamsFromBody(content)
	return ret
}

func (obj *KeyValue) ExtractParamsFromBody(body string) map[string]string {
	ret := make(map[string]string, 0)
	keyvalues := strings.Split(body, "&")
	for _, v := range keyvalues {
		kv := strings.Split(v, "=")
		if len(kv) == 2 {
			ret[kv[0]] = kv[1]
		}
	}
	return ret
}

//
//
//
type SendAccessTokenResult struct {
	KeyValue
}

func NewSendAccessTokenResult(content string) *SendAccessTokenResult {
	ret := new(SendAccessTokenResult)
	ret.KeyValues = ret.ExtractParamsFromBody(content)
	return ret
}

func (obj *SendAccessTokenResult) GetOAuthToken() string {
	return obj.KeyValues[OAuthToken]
}

func (obj *SendAccessTokenResult) GetOAuthTokenSecret() string {
	return obj.KeyValues[OAuthTokenSecret]
}

func (obj *SendAccessTokenResult) GetOAuthCallbackConfirmed() string {
	return obj.KeyValues[OAuthCallbackConfirmed]
}

//
func (obj *SendAccessTokenResult) GetUserID() string {
	return obj.KeyValues[UserID]
}

//
func (obj *SendAccessTokenResult) GetScreenName() string {
	return obj.KeyValues[ScreenName]
}

//
//
//
type SendRequestTokenResult struct {
	KeyValue
}

func NewSendRequestTokenResult(content string) *SendRequestTokenResult {
	ret := new(SendRequestTokenResult)
	ret.KeyValues = ret.ExtractParamsFromBody(content)
	return ret
}

func (obj *SendRequestTokenResult) GetOAuthToken() string {
	return obj.KeyValues[OAuthToken]
}

func (obj *SendRequestTokenResult) GetOAuthTokenSecret() string {
	return obj.KeyValues[OAuthTokenSecret]
}

func (obj *SendRequestTokenResult) GetOAuthCallbackConfirmed() string {
	return obj.KeyValues[OAuthCallbackConfirmed]
}

func (obj *SendRequestTokenResult) GetOAuthTokenUrl() string {
	return "https://api.twitter.com/oauth/authenticate?oauth_token=" + obj.GetOAuthToken()
}
