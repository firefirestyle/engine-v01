package blob

import (
	"encoding/json"

	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"

	"github.com/firefirestyle/engine-v01/prop"
)

type BlobManager struct {
	config BlobManagerConfig
}

type BlobManagerConfig struct {
	Kind        string
	PointerKind string
	CallbackUrl string
	HashLength  int
}

func NewBlobManager(config BlobManagerConfig) *BlobManager {
	ret := new(BlobManager)
	ret.config = config
	return ret
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

/**
 *
 *
 */
func (obj *BlobManager) SaveSignCache(ctx context.Context, dir, name, value string) {
	p := prop.NewMiniPath(dir)
	key := p.GetDir() + "/" + name
	if value == "" {
		memcache.Delete(ctx, key)
	} else {
		sk, _ := json.Marshal(map[string]interface{}{
			"k": obj.config.PointerKind,
			"n": key,
		})

		memcache.Set(ctx, &memcache.Item{
			Key:   string(sk),
			Value: []byte(value),
		})
	}
}

func (obj *BlobManager) LoadSignCache(ctx context.Context, dir, name string) (string, error) {
	p := prop.NewMiniPath(dir)
	key := p.GetDir() + "/" + name
	sk, _ := json.Marshal(map[string]interface{}{
		"k": obj.config.PointerKind,
		"n": key,
	})
	item, err := memcache.Get(ctx, string(sk))
	if err != nil {
		return "", err
	}

	return string(item.Value), nil
}
