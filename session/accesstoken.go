package session

import (
	"errors"
	"time"

	"encoding/json"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

var ErrorNotFound = errors.New("not found")
var ErrorAlreadyRegist = errors.New("already found")
var ErrorAlreadyUseMail = errors.New("already use mail")
var ErrorInvalid = errors.New("invalid")
var ErrorInvalidPass = errors.New("invalid password")
var ErrorOnServer = errors.New("server error")
var ErrorExtract = errors.New("failed to extract")

type GaeAccessTokenItem struct {
	UserName  string
	LoginTime time.Time

	LoginId   string `datastore:",noindex"`
	DeviceID  string `datastore:",noindex"`
	IP        string `datastore:",noindex"`
	Type      string `datastore:",noindex"`
	UserAgent string `datastore:",noindex"`
	Info      string `datastore:",noindex"`
}

type SessionManager struct {
	rootGroup          string
	MemcacheExpiration time.Duration
	loginIdKind        string
}

type AccessToken struct {
	gaeObject    *GaeAccessTokenItem
	gaeObjectKey *datastore.Key
	ItemKind     string
}

const (
	TypeRootGroup = "RootGroup"
	TypeUserName  = "UserName"
	TypeLoginTime = "LoginTime"
	TypeLoginId   = "LoginId"
	TypeDeviceID  = "DeviceID"
	TypeIP        = "IP"
	TypeType      = "Type"
	TypeInfo      = "Info"
	TypeUserAgent = "UserAgent"
)

func getStringFromProp(requestPropery map[string]interface{}, key string, defaultValue string) string {
	v := requestPropery[key]
	if v == nil {
		return defaultValue
	} else {
		return v.(string)
	}
}

//
// use package only. and in testcase
//
func (obj *AccessToken) ToJson() (string, error) {
	v := map[string]interface{}{
		TypeUserName:  obj.GetUserName(),             //
		TypeLoginTime: obj.GetLoginTime().UnixNano(), //
		TypeLoginId:   obj.GetLoginId(),              //
		TypeDeviceID:  obj.GetDeviceId(),             //
		TypeIP:        obj.GetIP(),                   //
		TypeType:      obj.gaeObject.Type,            //
		TypeUserAgent: obj.GetUserAgent(),            //
		TypeInfo:      obj.gaeObject.Info,
	}
	vv, e := json.Marshal(v)
	return string(vv), e
}

//
// use package only. and in testcase
//
func (obj *AccessToken) SetAccessTokenFromsJson(ctx context.Context, source string) error {
	v := make(map[string]interface{})
	e := json.Unmarshal([]byte(source), &v)
	if e != nil {
		return e
	}
	//
	obj.gaeObject.UserName = v[TypeUserName].(string)
	obj.gaeObject.LoginTime = time.Unix(0, int64(v[TypeLoginTime].(float64)))
	obj.gaeObject.LoginId = v[TypeLoginId].(string)

	obj.gaeObject.DeviceID = v[TypeDeviceID].(string)
	obj.gaeObject.IP = v[TypeIP].(string)
	obj.gaeObject.Type = v[TypeType].(string)
	obj.gaeObject.UserAgent = v[TypeUserAgent].(string)

	return nil
}

func (obj *AccessToken) GetLoginId() string {
	return obj.gaeObject.LoginId
}

func (obj *AccessToken) GetUserName() string {
	return obj.gaeObject.UserName
}

func (obj *AccessToken) GetIP() string {
	return obj.gaeObject.IP
}

func (obj *AccessToken) GetUserAgent() string {
	return obj.gaeObject.UserAgent
}

func (obj *AccessToken) GetDeviceId() string {
	return obj.gaeObject.DeviceID
}

func (obj *AccessToken) GetLoginTime() time.Time {
	return obj.gaeObject.LoginTime
}

func (obj *AccessToken) GetGAEObjectKey() *datastore.Key {
	return obj.gaeObjectKey
}

func (obj *AccessToken) LoadFromDB(ctx context.Context) error {
	source, errGet := memcache.Get(ctx, obj.gaeObjectKey.StringID())
	if errGet == nil {
		errSet := obj.SetAccessTokenFromsJson(ctx, string(source.Value))
		if errSet == nil {
			return nil
		} else {
		}
	}

	errGetFromStore := datastore.Get(ctx, obj.gaeObjectKey, obj.gaeObject)
	if errGetFromStore == nil {
		obj.UpdateMemcache(ctx)
	}
	return errGetFromStore
}

func (obj *AccessToken) Logout(ctx context.Context) error {
	obj.gaeObject.LoginId = ""
	return obj.Save(ctx)
}

func (obj *AccessToken) Save(ctx context.Context) error {
	_, e := datastore.Put(ctx, obj.gaeObjectKey, obj.gaeObject)
	if e != nil {
		obj.UpdateMemcache(ctx)
	}
	return e
}

func (obj *AccessToken) DeleteFromDB(ctx context.Context) error {
	memcache.Delete(ctx, obj.gaeObjectKey.StringID())
	return datastore.Delete(ctx, obj.gaeObjectKey)
}

func (obj *AccessToken) UpdateMemcache(ctx context.Context) error {
	userObjMemSource, err_toJson := obj.ToJson()
	if err_toJson == nil {
		userObjMem := &memcache.Item{
			Key:   obj.gaeObjectKey.StringID(),
			Value: []byte(userObjMemSource), //
		}
		memcache.Set(ctx, userObjMem)
	}
	return err_toJson
}
