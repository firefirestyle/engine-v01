package article

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"crypto/rand"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"

	"encoding/json"

	"google.golang.org/appengine/memcache"
)

type ArticleManagerConfig struct {
	KindArticle    string
	KindPointer    string
	LimitOfFinding int
	LengthHash     int
}

type ArticleManager struct {
	config ArticleManagerConfig
}

func NewArticleManager(config ArticleManagerConfig) *ArticleManager {
	if config.KindArticle == "" {
		config.KindArticle = "fa"
	}
	if config.KindPointer == "" {
		config.KindPointer = config.KindArticle + "-pointer"
	}
	if config.LimitOfFinding <= 0 {
		config.LimitOfFinding = 20
	}
	ret := new(ArticleManager)
	ret.config = config
	return ret
}

func (obj *ArticleManager) GetKind() string {
	return obj.config.KindArticle
}

func (obj *ArticleManager) MakeArticleId(created time.Time, secretKey string) string {
	hashKey := obj.hashStr(fmt.Sprintf("s:%s;c:%d;", secretKey, created.UnixNano()))
	return hashKey
}

func (obj *ArticleManager) MakeStringId(articleId string, sign string) string {
	propObj := miniprop.NewMiniProp()
	propObj.SetString("i", articleId)
	propObj.SetString("s", sign)
	return string(propObj.ToJson())
}

type StringIdInfo struct {
	ArticleId string
	Sign      string
}

func (obj *ArticleManager) ExtractInfoFromStringId(stringId string) *StringIdInfo {
	propObj := miniprop.NewMiniPropFromJson([]byte(stringId))
	return &StringIdInfo{
		ArticleId: propObj.GetString("i", ""),
		Sign:      propObj.GetString("s", ""),
	}
}

func (obj *ArticleManager) SaveArticleWithImmutable(ctx context.Context, artObj *Article) (*Article, error) {
	sign := strconv.Itoa(time.Now().Nanosecond())
	nextArObj := obj.NewArticleFromArticle(ctx, artObj, sign)
	nextArObj.SetUpdated(time.Now())
	saveErr := nextArObj.saveOnDB(ctx)
	if saveErr != nil {
		return artObj, saveErr
	}

	if artObj.gaeObject.Sign != "0" {
		obj.DeleteFromArticleId(ctx, artObj.GetArticleId(), artObj.GetSign())
	}
	obj.SaveSignCache(ctx, nextArObj.GetArticleId(), sign)
	return nextArObj, nil
}

func (obj *ArticleManager) GetLimitOfFinding() int {
	return obj.config.LimitOfFinding
}

//
//
//

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

func (obj *ArticleManager) hash(v string) string {
	sha1Obj := sha1.New()
	sha1Obj.Write([]byte(v))
	return string(sha1Obj.Sum(nil))
}

func (obj *ArticleManager) hashStr(v string) string {
	sha1Obj := sha1.New()
	sha1Obj.Write([]byte(v))
	articleIdHash := string(base32.StdEncoding.EncodeToString(sha1Obj.Sum(nil)))
	if obj.config.LengthHash > 5 && len(articleIdHash) > obj.config.LengthHash {
		articleIdHash = articleIdHash[:obj.config.LengthHash]
	}
	return articleIdHash
}

func (obj *ArticleManager) makeRandomId() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

/**
 *
 *
 */
func (obj *ArticleManager) SaveSignCache(ctx context.Context, key, value string) {
	if value == "" {
		memcache.Delete(ctx, key)
	} else {
		sk, _ := json.Marshal(map[string]interface{}{
			"k": obj.config.KindPointer,
			"n": key,
		})

		memcache.Set(ctx, &memcache.Item{
			Key:   string(sk),
			Value: []byte(value),
		})
	}
}

func (obj *ArticleManager) LoadSignCache(ctx context.Context, key string) (string, error) {
	sk, _ := json.Marshal(map[string]interface{}{
		"k": obj.config.KindPointer,
		"n": key,
	})
	item, err := memcache.Get(ctx, string(sk))
	if err != nil {
		return "", err
	}

	return string(item.Value), nil
}
