package article

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine/memcache"
)

type GaeObjectArticle struct {
	UserName   string
	Title      string   `datastore:",noindex"`
	Tags       []string `datastore:"Tags.Tag"`
	Point      float64
	Lat        float64
	Lng        float64
	PropNames  []string `datastore:"Props.Name"`
	PropValues []string `datastore:"Props.Value"`
	Cont       string   `datastore:",noindex"`
	Info       string   `datastore:",noindex"`
	Sign       string   `datastore:",noindex"`
	ArticleId  string
	Created    time.Time
	Updated    time.Time
	SecretKey  string `datastore:",noindex"`
	IconUrl    string `datastore:",noindex"`
}

type Article struct {
	gaeObjectKey *datastore.Key
	gaeObject    *GaeObjectArticle
	kind         string
}

const (
	TypeUserName   = "UserName"
	TypeTitle      = "Title"
	TypeTag        = "Tag"
	TypePoint      = "Point"
	TypePropNames  = "PropNames"
	TypePropValues = "PropValues"
	TypeCont       = "Cont"
	TypeInfo       = "Info"
	TypeType       = "Type"
	TypeSign       = "Sign"
	TypeArticleId  = "ArticleId"
	TypeCreated    = "Created"
	TypeUpdated    = "Updated"
	TypeSecretKey  = "SecretKey"
	TypeTarget     = "Target"
	TypeLat        = "Lat"
	TypeLng        = "Lng"
	TypeIconUrl    = "IconUrl"
)

func (obj *Article) updateMemcache(ctx context.Context) error {
	userObjMemSource := obj.ToJson()
	userObjMem := &memcache.Item{
		Key:   obj.gaeObjectKey.StringID(),
		Value: []byte(userObjMemSource), //
	}
	memcache.Set(ctx, userObjMem)
	return nil
}

//
//
//
func (obj *Article) GetGaeObjectKind() string {
	return obj.kind
}

func (obj *Article) GetGaeObjectKey() *datastore.Key {
	return obj.gaeObjectKey
}

func (obj *Article) GetUserName() string {
	return obj.gaeObject.UserName
}

func (obj *Article) GetSign() string {
	return obj.gaeObject.Sign
}

func (obj *Article) GetIconUrl() string {
	return obj.gaeObject.IconUrl
}

func (obj *Article) SetIconUrl(v string) {
	obj.gaeObject.IconUrl = v
}

func (obj *Article) GetInfo() string {
	return obj.gaeObject.Info
}

func (obj *Article) SetInfo(v string) {
	obj.gaeObject.Info = v
}

func (obj *Article) SetUserName(v string) {
	obj.gaeObject.UserName = v
}

func (obj *Article) GetTitle() string {
	return obj.gaeObject.Title
}

func (obj *Article) SetTitle(v string) {
	obj.gaeObject.Title = v
}

func (obj *Article) GetTags() []string {
	ret := make([]string, 0)
	for _, v := range obj.gaeObject.Tags {
		ret = append(ret, v)
	}
	return ret
}

func (obj *Article) SetTags(vs []string) {
	obj.gaeObject.Tags = make([]string, 0)
	for _, v := range vs {
		obj.gaeObject.Tags = append(obj.gaeObject.Tags, v)
	}
}

func (obj *Article) GetCont() string {
	return obj.gaeObject.Cont
}

func (obj *Article) SetCont(v string) {
	obj.gaeObject.Cont = v
}

func (obj *Article) GetParentId() string {
	return obj.gaeObject.Sign
}

func (obj *Article) SetParentId(v string) {
	obj.gaeObject.Sign = v
}

func (obj *Article) GetArticleId() string {
	return obj.gaeObject.ArticleId
}

func (obj *Article) GetCreated() time.Time {
	return obj.gaeObject.Created
}

func (obj *Article) GetUpdated() time.Time {
	return obj.gaeObject.Updated
}

func (obj *Article) SetUpdated(v time.Time) {
	obj.gaeObject.Updated = v
}

func (obj *Article) GetPoint() float64 {
	return obj.gaeObject.Point
}

func (obj *Article) SetPoint(v float64) {
	obj.gaeObject.Point = v
}

func (obj *Article) GetLat() float64 {
	return obj.gaeObject.Lat
}

func (obj *Article) SetLat(v float64) {
	obj.gaeObject.Lat = v
}

func (obj *Article) GetLng() float64 {
	return obj.gaeObject.Lng
}

func (obj *Article) SetLng(v float64) {
	obj.gaeObject.Lng = v
}

func (obj *Article) ClearProp() {
	obj.gaeObject.PropNames = nil
	obj.gaeObject.PropValues = nil
}

func (obj *Article) GetProp(name string) string {
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
	p := miniprop.NewMiniPropFromJson([]byte(obj.gaeObject.PropValues[index]))
	return p.GetString(name, "")
}

func (obj *Article) SetProp(name, v string) {
	index := -1
	p := miniprop.NewMiniProp()
	p.SetString(name, v)
	v = string(p.ToJson())
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
