package twitter

//
// https://dev.twitter.com/oauth/overview
// https://dev.twitter.com/web/sign-in/implementing
//

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	//	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"
	//"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

const (
	OAuth1Callback                = "oauth_callback"
	OAuth1ConsumerKey             = "oauth_consumer_key"
	OAuth1SignatureMethod         = "oauth_signature_method"
	OAuth1SignatureMethodHmacSHA1 = "HMAC-SHA1"
	OAuth1Version                 = "oauth_version"
	OAuth1Version1                = "1.0"
	OAuth1Nonce                   = "oauth_nonce"
	OAuth1TIme                    = "oauth_timestamp"
	OAuth1Signature               = "oauth_signature"
	OAuth1Token                   = "oauth_token"
	OAuthAuthorizationHeader      = "Authorization"
)

type OAuth1Client struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	Callback          string
	Method            string
	Version           string
	AuthParam         map[string]string
	AllowInvalidSSL   bool
}

func NewOAuthClient(consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string, allowInvalidSSL bool) *OAuth1Client {
	ret := new(OAuth1Client)
	ret.ConsumerKey = consumerKey
	ret.ConsumerSecret = consumerSecret
	ret.AccessToken = accessToken
	ret.AccessTokenSecret = accessTokenSecret
	ret.Method = OAuth1SignatureMethodHmacSHA1
	ret.Version = OAuth1Version1
	ret.AuthParam = make(map[string]string, 0)
	ret.AllowInvalidSSL = allowInvalidSSL
	return ret
}

func (obj *OAuth1Client) Post(ctx context.Context,
	urlStr string, headers map[string]string, body string) (string, error) {

	obj.Clear(urlStr)
	obj.Sign(urlStr)

	request, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	//
	v := obj.MakeAuthorizationHeader()
	request.Header.Add(OAuthAuthorizationHeader, v)

	//
	//
	var client *http.Client
	if obj.AllowInvalidSSL == true {
		client = &http.Client{
			Transport: &urlfetch.Transport{
				Context: ctx,
				AllowInvalidServerCertificate: true,
				//			Deadline: 10 * time.Second,
			},
		}
	} else {
		client = urlfetch.Client(ctx)
	}

	//	client.Transport.

	response, err1 := client.Do(request)
	if err1 != nil {
		return "", err1
	}
	//
	result := make([]byte, 256)
	_, err = response.Body.Read(result)
	if err != nil {
		return "", err
	}
	//
	return string(result), nil
}

func (obj *OAuth1Client) MakeNonce() string {
	var n, m uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	binary.Read(rand.Reader, binary.LittleEndian, &m)
	return strconv.FormatUint(n, 36) + strconv.FormatUint(m, 36)
}

func (obj *OAuth1Client) MakeTimestamp() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func (obj *OAuth1Client) Clear(targetAddr string) {
	obj.AuthParam = make(map[string]string, 0)
	if obj.Callback != "" {
		obj.AuthParam[OAuth1Callback] = obj.Callback
	}
	obj.AuthParam[OAuth1ConsumerKey] = obj.ConsumerKey
	obj.AuthParam[OAuth1SignatureMethod] = obj.Method
	obj.AuthParam[OAuth1Version] = obj.Version
	//
	obj.AuthParam[OAuth1Nonce] = obj.MakeNonce()
	if obj.AccessToken != "" {
		obj.AuthParam[OAuth1Token] = obj.AccessToken
	}

	obj.AuthParam[OAuth1TIme] = obj.MakeTimestamp()
}

func (obj *OAuth1Client) Sign(targetAddr string) {
	obj.AuthParam[OAuth1Signature] = url.QueryEscape(obj.MakeSignature(obj.MakeSignBaseString(targetAddr)))
}

func (obj *OAuth1Client) MakeSignBaseString(targetAddr string) string {
	ret := "POST&" + url.QueryEscape(targetAddr) + "&"
	//
	keys := make([]string, 0)
	for key, _ := range obj.AuthParam {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	//
	params := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		params[i] = fmt.Sprintf("%s=%s",
			url.QueryEscape(keys[i]), //
			url.QueryEscape(obj.AuthParam[keys[i]]))
	}
	paramSignSrc := strings.Join(params, "&")
	//
	ret += url.QueryEscape(paramSignSrc)
	return ret
}

func (obj *OAuth1Client) MakeSignature(signBase string) string {
	signKey := url.QueryEscape(obj.ConsumerSecret) + "&" + url.QueryEscape(obj.AccessTokenSecret)
	hmacObj := hmac.New(sha1.New, []byte(signKey))
	hmacObj.Write([]byte(signBase))
	retSrc := hmacObj.Sum(nil)
	return base64.StdEncoding.EncodeToString(retSrc)
}

func (obj *OAuth1Client) MakeAuthorizationHeader() string {
	params := make([]string, 0)
	for k, v := range obj.AuthParam {
		if k == OAuth1Callback {
			params = append(params, fmt.Sprintf(`%s="%s"`, k, url.QueryEscape(v)))
		} else {
			params = append(params, fmt.Sprintf(`%s="%s"`, k, v))
		}
	}
	return fmt.Sprintf("OAuth %s", strings.Join(params, ","))
}
