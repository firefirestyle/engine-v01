package session

import (
	"time"

	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/mssola/user_agent"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type SessionManagerConfig struct {
	Kind string
}

func NewSessionManager(config SessionManagerConfig) *SessionManager {
	ret := new(SessionManager)
	if config.Kind == "" {
		ret.loginIdKind = "LoginId"
	} else {
		ret.loginIdKind = config.Kind
	}
	return ret
}

func (obj *SessionManager) NewAccessToken(ctx context.Context, userName string, config OptionInfo) *AccessToken {
	ret := new(AccessToken)
	ret.gaeObject = new(GaeAccessTokenItem)
	loginTime := time.Now()
	idInfoObj := obj.MakeLoginIdInfo(userName, config)

	ret.gaeObject.LoginId = idInfoObj.LoginId
	ret.gaeObject.IP = config.IP
	ret.gaeObject.Type = config.LoginType
	ret.gaeObject.LoginTime = loginTime
	ret.gaeObject.DeviceID = idInfoObj.DeviceId
	ret.gaeObject.UserName = userName
	ret.gaeObject.UserAgent = config.UserAgent

	ret.ItemKind = obj.loginIdKind
	ret.gaeObjectKey = obj.NewAccessTokenGaeObjectKey(ctx, idInfoObj)

	return ret
}

func (obj *SessionManager) LoadAccessToken(ctx context.Context, loginId string) (*AccessToken, error) {
	idInfo, err := obj.MakeSourceFromAccessToken(loginId)
	if err != nil {
		return nil, err
	}
	accessTokenObj := new(AccessToken)
	accessTokenObj.ItemKind = obj.loginIdKind
	accessTokenObj.gaeObject = new(GaeAccessTokenItem)
	accessTokenObj.gaeObjectKey = obj.NewAccessTokenGaeObjectKey(ctx, idInfo)
	accessTokenObj.gaeObject.LoginId = loginId

	err = accessTokenObj.LoadFromDB(ctx)
	if err != nil {
		return nil, err
	}
	return accessTokenObj, nil
}

func (obj *SessionManager) NewAccessTokenGaeObjectKey(ctx context.Context, idInfoObj Source) *datastore.Key {
	return datastore.NewKey(ctx, obj.loginIdKind, obj.MakeGaeObjectKeyStringId(idInfoObj.UserName, idInfoObj.DeviceId), 0, nil)
}

func (obj *SessionManager) MakeGaeObjectKeyStringId(userName string, deviceId string) string {
	return obj.loginIdKind + ":" + obj.rootGroup + ":" + userName + ":" + deviceId
}

func (obj *SessionManager) MakeSourceFromAccessToken(loginId string) (Source, error) {
	binary := []byte(loginId)
	if len(binary) <= 28+28+1 {
		return Source{}, ErrorExtract
	}
	//
	binaryUser, err := base64.StdEncoding.DecodeString(string(binary[28*2:]))
	if err != nil {
		return Source{}, ErrorExtract
	}
	//
	return Source{
		DeviceId: string(binary[28 : 28*2]),
		UserName: string(binaryUser),
	}, nil
}

func (obj *SessionManager) MakeDeviceId(userName string, info OptionInfo) string {
	uaObj := user_agent.New(info.UserAgent)
	sha1Hash := sha1.New()
	b, _ := uaObj.Browser()
	io.WriteString(sha1Hash, b)
	io.WriteString(sha1Hash, uaObj.OS())
	io.WriteString(sha1Hash, uaObj.Platform())
	return base64.StdEncoding.EncodeToString(sha1Hash.Sum(nil))
}

func (obj *SessionManager) MakeLoginIdInfo(userName string, config OptionInfo) Source {
	deviceID := obj.MakeDeviceId(userName, config)
	loginId := ""
	sha1Hash := sha1.New()
	io.WriteString(sha1Hash, deviceID)
	io.WriteString(sha1Hash, userName)
	io.WriteString(sha1Hash, fmt.Sprintf("%X", rand.Int63()))
	loginId = base64.StdEncoding.EncodeToString(sha1Hash.Sum(nil))
	loginId += deviceID
	loginId += base64.StdEncoding.EncodeToString([]byte(userName))
	return Source{
		DeviceId: deviceID,
		UserName: userName,
		LoginId:  loginId,
	}
}

func (obj *SessionManager) CheckAccessToken(ctx context.Context, accessToken string, option OptionInfo, isIPCheck bool) CheckResult {
	accessTokenObj, err := obj.LoadAccessToken(ctx, accessToken)
	if err != nil {
		return CheckResult{
			IsLogin:        false,
			AccessTokenObj: nil,
		}
	}

	// todos
	if accessTokenObj.GetLoginId() != accessToken {
		return CheckResult{
			IsLogin:        false,
			AccessTokenObj: accessTokenObj,
		}
	}

	//
	if isIPCheck == true {
		if accessTokenObj.GetDeviceId() != obj.MakeDeviceId(accessTokenObj.GetUserName(), option) {
			return CheckResult{
				IsLogin:        false,
				AccessTokenObj: accessTokenObj,
			}
		}
	}

	return CheckResult{
		IsLogin:        true,
		AccessTokenObj: accessTokenObj,
	}
}

func (obj *SessionManager) Login(ctx context.Context, userName string, option OptionInfo) (*AccessToken, error) {
	accessTokenObj := obj.NewAccessToken(ctx, userName, option)
	err1 := accessTokenObj.Save(ctx)
	return accessTokenObj, err1
}

func (obj *SessionManager) Logout(ctx context.Context, accessToken string, option OptionInfo) error {
	checkResult := obj.CheckAccessToken(ctx, accessToken, option, false)
	if checkResult.IsLogin == false {
		return nil
	}
	return checkResult.AccessTokenObj.Logout(ctx)
}

//
//
//
type OptionInfo struct {
	IP        string
	UserAgent string
	LoginType string
}

func MakeOptionInfo(r *http.Request) OptionInfo {
	return OptionInfo{IP: r.RemoteAddr, UserAgent: r.UserAgent()}
}

//
//
//
type Source struct {
	DeviceId string
	UserName string
	LoginId  string
}

type CheckResult struct {
	IsLogin        bool
	AccessTokenObj *AccessToken
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
