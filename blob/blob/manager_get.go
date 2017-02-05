package blob

import (
	"golang.org/x/net/context"

	"errors"

	"google.golang.org/appengine/datastore"
)

func (obj *BlobManager) GetBlobItemFromGaeKey(ctx context.Context, gaeKey *datastore.Key) (*BlobItem, error) {
	memCacheObj, errMemCcache := obj.NewBlobItemFromMemcache(ctx, gaeKey.StringID())
	if errMemCcache == nil {
		return memCacheObj, nil
	}
	//
	//
	var item GaeObjectBlobItem
	err := datastore.Get(ctx, gaeKey, &item)
	if err != nil {
		return nil, err
	}
	ret := new(BlobItem)
	ret.gaeObject = &item
	ret.gaeKey = gaeKey

	if err == nil {
		ret.updateMemcache(ctx)
	}
	return ret, nil
}

func (obj *BlobManager) GetBlobItem(ctx context.Context, parent string, name string, sign string) (*BlobItem, error) {

	key := obj.NewBlobItemGaeKey(ctx, parent, name, sign)

	return obj.GetBlobItemFromGaeKey(ctx, key)
}

func (obj *BlobManager) GetBlobItemStringIdFromQuery(ctx context.Context, parent string, name string) (string, error) {
	founded := obj.FindBlobItemFromPath(ctx, parent, name, "")
	if len(founded.Keys) <= 0 {
		return "", errors.New("not found blobitem")
	}
	key := obj.NewBlobItemGaeKeyFromStringId(ctx, founded.Keys[0])
	return key.StringID(), nil
}

func (obj *BlobManager) GetBlobItemFromQuery(ctx context.Context, parent string, name string) (*BlobItem, error) {
	founded := obj.FindBlobItemFromPath(ctx, parent, name, "")
	if len(founded.Keys) <= 0 {
		return nil, errors.New("not found blobitem")
	}
	key := obj.NewBlobItemGaeKeyFromStringId(ctx, founded.Keys[0])
	return obj.GetBlobItemFromGaeKey(ctx, key)
}

func (obj *BlobManager) GetBlobItemFromStringId(ctx context.Context, stringId string) (*BlobItem, error) {
	key := obj.NewBlobItemGaeKeyFromStringId(ctx, stringId)
	return obj.GetBlobItemFromGaeKey(ctx, key)
}

func (obj *BlobManager) GetBlobItemFromPointer(ctx context.Context, parent string, name string) (*BlobItem, error) {
	sign, err := obj.LoadSignCache(ctx, parent, name)
	if err != nil {
		b, e := obj.GetBlobItem(ctx, parent, name, sign)
		if e == nil {
			return b, e
		}
	}
	o, e := obj.GetBlobItemFromQuery(ctx, parent, name)
	return o, e
}

func (obj *BlobManager) GetBlobItemStringIdFromPointer(ctx context.Context, parent string, name string) (string, string, error) {
	sign, err := obj.LoadSignCache(ctx, parent, name)
	if err != nil {
		return obj.MakeStringId(parent, name, sign), obj.MakeBlobId(parent, name), nil
	} else {
		o, e := obj.GetBlobItemStringIdFromQuery(ctx, parent, name)
		return o, obj.MakeBlobId(parent, name), e
	}
}
