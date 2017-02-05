package user

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

//
// save and delete
//
func (obj *User) updateMemcache(ctx context.Context) error {
	userObjMemSource := obj.ToJson()
	userObjMem := &memcache.Item{
		Key:   obj.gaeObjectKey.StringID(),
		Value: []byte(userObjMemSource), //
	}
	return memcache.Set(ctx, userObjMem)
}

func (obj *User) deleteMemcache(ctx context.Context) error {
	return memcache.Delete(ctx, obj.gaeObjectKey.StringID())
}

//
//
//
func (obj *User) loadFromDB(ctx context.Context) error {
	item, err := memcache.Get(ctx, obj.gaeObjectKey.StringID())
	if err == nil {
		err1 := obj.SetUserFromsJson(ctx, string(item.Value))
		if err1 == nil {
			return nil
		} else {
			log.Infof(ctx, ">>> Failed Load UseObj Json On LoadFronDB :: %s", err1.Error())
		}
	}
	//
	//
	err_loaded := datastore.Get(ctx, obj.gaeObjectKey, obj.gaeObject)
	if err_loaded != nil {
		return err_loaded
	}

	obj.updateMemcache(ctx)

	return nil
}

func (obj *User) pushToDB(ctx context.Context) error {
	_, e := datastore.Put(ctx, obj.gaeObjectKey, obj.gaeObject)
	obj.updateMemcache(ctx)
	return e
}
