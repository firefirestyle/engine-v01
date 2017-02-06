package blob

import (
	"golang.org/x/net/context"

	"time"

	"github.com/firefirestyle/engine-v01/prop"
	m "github.com/firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
	"google.golang.org/appengine/blobstore"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

func (obj *BlobManager) NewBlobItem(ctx context.Context, parent string, name string, blobKey string) *BlobItem {
	ret := new(BlobItem)
	ret.gaeObject = new(GaeObjectBlobItem)
	{
		p := prop.NewMiniPath(parent)
		ret.gaeObject.Parent = p.GetDir()
	}
	ret.gaeObject.Name = name
	ret.gaeObject.BlobKey = blobKey
	ret.gaeObject.Updated = time.Now()
	ret.gaeObject.Sign = blobKey
	ret.gaeKey = datastore.NewKey(ctx, obj.config.Kind, obj.MakeStringId(parent, name, blobKey), 0, nil)
	return ret
}

func (obj *BlobManager) NewBlobItemFromMemcache(ctx context.Context, keyId string) (*BlobItem, error) {
	jsonSource, errGetJsonSource := memcache.Get(ctx, keyId)
	if errGetJsonSource != nil {
		return nil, errGetJsonSource
	}

	ret := new(BlobItem)
	ret.gaeKey = datastore.NewKey(ctx, obj.config.Kind, keyId, 0, nil)
	ret.gaeObject = new(GaeObjectBlobItem)
	err := ret.SetParamFromJson(jsonSource.Value)
	return ret, err
}

func (obj *BlobManager) NewBlobItemGaeKey(ctx context.Context, parent string, name string, sign string) *datastore.Key {
	return obj.NewBlobItemGaeKeyFromStringId(ctx, obj.MakeStringId(parent, name, sign))
}

func (obj *BlobManager) NewBlobItemGaeKeyFromStringId(ctx context.Context, stringId string) *datastore.Key {
	return datastore.NewKey(ctx, obj.config.Kind, stringId, 0, nil)
}

func (obj *BlobItem) updateMemcache(ctx context.Context) error {
	userObjMemSource, err_toJson := obj.ToJson()
	if err_toJson == nil {
		userObjMem := &memcache.Item{
			Key:   obj.gaeKey.StringID(),
			Value: []byte(userObjMemSource), //
		}
		memcache.Set(ctx, userObjMem)
	}
	return err_toJson
}

func (obj *BlobItem) saveDB(ctx context.Context) error {
	_, e := datastore.Put(ctx, obj.gaeKey, obj.gaeObject)
	obj.updateMemcache(ctx)
	return e
}

func (obj *BlobManager) DeleteBlobItemFromStringId(ctx context.Context, stringId string) error {
	keyInfo := obj.GetKeyInfoFromStringId(stringId)
	blobKey := keyInfo.Sign
	if blobKey != "" {
		if nil != blobstore.Delete(ctx, appengine.BlobKey(blobKey)) {
			Debug(ctx, "GOMIDATA in DeleteFromDBFromStringId : "+stringId)
		}
	}
	return datastore.Delete(ctx, obj.NewBlobItemGaeKeyFromStringId(ctx, stringId))
}

type BlobItemKeyInfo struct {
	Kind   string
	Parent string
	Name   string
	Sign   string
}

func (obj *BlobManager) GetKeyInfoFromStringId(stringId string) BlobItemKeyInfo {
	propObj := m.NewMiniPropFromJson([]byte(stringId))
	return BlobItemKeyInfo{
		Kind:   propObj.GetString("k", ""),
		Parent: propObj.GetString("d", ""),
		Name:   propObj.GetString("f", ""),
		Sign:   propObj.GetString("s", ""),
	}
}

func (obj *BlobManager) MakeStringId(parent string, name string, sign string) string {
	propObj := m.NewMiniProp()
	propObj.SetString("k", obj.config.Kind)
	propObj.SetString("d", parent)
	propObj.SetString("f", name)
	propObj.SetString("s", sign)
	return string(propObj.ToJson())
}

func (obj *BlobManager) MakeBlobId(parent string, name string) string {
	propObj := m.NewMiniProp()
	propObj.SetString("d", parent)
	propObj.SetString("f", name)
	return string(propObj.ToJson())
}
