package user

import (
	"crypto/sha1"
	"io"
	"time"

	"encoding/base32"

	"strconv"

	"github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

const (
	UserStatePublic  = "public"
	UserStatePrivate = "private"
	UserStateAll     = ""
)

const (
	TypeRootGroup   = "RootGroup"
	TypeDisplayName = "DisplayName"
	TypeUserName    = "UserName"
	TypeCreated     = "Created"
	TypeUpdated     = "Updated"
	TypeState       = "State"
	TypeTag         = "Tag"
	TypePublicInfo  = "PublicInfo"
	TypePoint       = "Point"
	TypePropNames   = "PropNames"
	TypePropValues  = "PropValues"
	TypeIconUrl     = "IconUrl"
	TypePrivateInfo = "PrivateInfo"
	TypeSign        = "Sign"
	TypeCont        = "Cont"
	TypePermission  = "Permission"
)

type GaeUserItem struct {
	DisplayName string
	UserName    string
	Created     time.Time
	Updated     time.Time
	State       string
	PublicInfo  string   `datastore:",noindex"`
	PrivateInfo string   `datastore:",noindex"`
	Tags        []string `datastore:"Tags.Tag"`
	PropNames   []string `datastore:"Props.Name"`
	PropValues  []string `datastore:"Props.Value"`
	Point       float64
	IconUrl     string `datastore:",noindex"`
	Sign        string `datastore:",noindex"`
	Cont        string `datastore:",noindex"`
	Permission  int
}

type User struct {
	gaeObject    *GaeUserItem
	gaeObjectKey *datastore.Key
	kind         string
	prop         map[string]map[string]interface{}
}

// ----
// new object
// ----

func (obj *UserManager) newUserGaeObjectKey(ctx context.Context, userName, sign string) *datastore.Key {
	return datastore.NewKey(ctx, obj.config.UserKind, obj.MakeStringId(userName, sign), 0, nil)
}

func (obj *UserManager) newUserWithUserName(ctx context.Context) *User {
	var userObj *User = nil
	var err error = nil
	for {
		hashObj := sha1.New()
		now := time.Now().UnixNano()
		io.WriteString(hashObj, prop.MakeRandomId())
		io.WriteString(hashObj, strconv.FormatInt(now, 36))
		userName := string(base32.StdEncoding.EncodeToString(hashObj.Sum(nil)))
		if obj.config.LengthHash >= 5 && len(userName) > obj.config.LengthHash {
			userName = userName[:obj.config.LengthHash]
		}
		userObj, err = obj.newUser(ctx, userName, strconv.FormatInt(now, 16))
		if err != nil {
			break
		}
	}
	return userObj
}

func (obj *UserManager) newUser(ctx context.Context, userName string, sign string) (*User, error) {
	ret := new(User)
	ret.prop = make(map[string]map[string]interface{})
	ret.kind = obj.config.UserKind
	ret.gaeObject = new(GaeUserItem)
	ret.gaeObject.UserName = userName
	ret.gaeObjectKey = obj.newUserGaeObjectKey(ctx, userName, sign)
	e := ret.loadFromDB(ctx)
	return ret, e
}

func (obj *UserManager) cloneUser(ctx context.Context, user *User) *User {
	ret := new(User)
	now := time.Now().UnixNano()
	ret.gaeObject = new(GaeUserItem)
	ret.SetUserFromsMap(ctx, user.ToMapAll())
	ret.gaeObjectKey = obj.newUserGaeObjectKey(ctx, ret.GetUserName(), strconv.FormatInt(now, 36))
	return ret
}

func (obj *UserManager) newUserFromStringID(ctx context.Context, stringId string) *User {
	ret := new(User)
	ret.prop = make(map[string]map[string]interface{})
	ret.kind = obj.config.UserKind
	ret.gaeObject = new(GaeUserItem)
	ret.gaeObjectKey = datastore.NewKey(ctx, obj.config.UserKind, stringId, 0, nil)
	return ret
}

// ----
// getter setter
// ----

func (obj *User) GetUserName() string {
	return obj.gaeObject.UserName
}

func (obj *User) GetDisplayName() string {
	return obj.gaeObject.DisplayName
}

func (obj *User) SetDisplayName(v string) {
	obj.gaeObject.DisplayName = v
}

func (obj *User) GetHaveIcon() bool {
	if obj.gaeObject.IconUrl == "" {
		return false
	} else {
		return true
	}
}

func (obj *User) SetIconUrl(v string) {
	obj.gaeObject.IconUrl = v
}

func (obj *User) GetIconUrl() string {
	return obj.gaeObject.IconUrl
}

func (obj *User) GetCreated() time.Time {
	return obj.gaeObject.Created
}

func (obj *User) GetLogined() time.Time {
	return obj.gaeObject.Updated
}

func (obj *User) SetLogined(v time.Time) {
	obj.gaeObject.Updated = v
}

func (obj *User) GetPublicInfo() string {
	return obj.gaeObject.PublicInfo
}

func (obj *User) SetPublicInfo(v string) {
	obj.gaeObject.PublicInfo = v
}

func (obj *User) GetPrivateInfo() string {
	return obj.gaeObject.PrivateInfo
}

func (obj *User) SetPrivateInfo(v string) {
	obj.gaeObject.PrivateInfo = v
}

func (obj *User) GetPoint() float64 {
	return obj.gaeObject.Point
}

func (obj *User) SetPoint(v float64) {
	obj.gaeObject.Point = v
}

func (obj *User) GetProp(name string) string {
	index := -1
	for i, v := range obj.gaeObject.PropNames {
		if v == name {
			index = i
			break
		}
	}
	if index < 0 {
		return ""
	}
	p := prop.NewMiniPropFromJson([]byte(obj.gaeObject.PropValues[index]))
	return p.GetString(name, "")
}

func (obj *User) SetProp(name, v string) {
	index := -1
	v = MakePropValue(name, v)
	for i, iv := range obj.gaeObject.PropNames {
		if iv == name {
			index = i
			break
		}
	}
	if index == -1 {
		obj.gaeObject.PropValues = append(obj.gaeObject.PropValues, v)
		obj.gaeObject.PropNames = append(obj.gaeObject.PropNames, name)
	} else {
		obj.gaeObject.PropValues[index] = v
	}
}

func (obj *User) SetStatus(v string) {
	obj.gaeObject.State = v
}

func (obj *User) GetStatus() string {
	return obj.gaeObject.State
}

func (obj *User) SetCont(v string) {
	obj.gaeObject.Cont = v
}

func (obj *User) GetCont() string {
	return obj.gaeObject.Cont
}

func (obj *User) GetSign() string {
	return obj.gaeObject.Sign
}

func (obj *User) GetStringId() string {
	return obj.gaeObjectKey.StringID()
}

func (obj *User) GetTags() []string {
	ret := make([]string, 0)
	for _, v := range obj.gaeObject.Tags {
		ret = append(ret, v)
	}
	return ret
}

func (obj *User) SetPermission(v int) {
	obj.gaeObject.Permission = v
}

func (obj *User) GetPermission() int {
	return obj.gaeObject.Permission
}

func (obj *User) IsMaster() bool {
	return obj.gaeObject.Permission&0x80 == 0x80
}

func (obj *User) SetMaster(on bool) {
	if on == true {
		obj.gaeObject.Permission = obj.gaeObject.Permission | 0x00000080
	} else {
		obj.gaeObject.Permission = obj.gaeObject.Permission & 0xFFFFFF7F
	}
}

func (obj *User) SetTags(vs []string) {
	obj.gaeObject.Tags = make([]string, 0)
	for _, v := range vs {
		obj.gaeObject.Tags = append(obj.gaeObject.Tags, v)
	}
}
