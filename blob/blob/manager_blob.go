package blob

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"net/http"
	"net/url"

	"golang.org/x/net/context"

	"io"

	"encoding/base32"
	"encoding/base64"
	"strings"

	"sort"

	"github.com/kyorohiro/ramenhunter_sv_lib/go.miniprop"
	"google.golang.org/appengine/blobstore"
)

//
// for make original hundler
//

func (obj *BlobManager) MakeRequestUrl(ctx context.Context, dirName string, fileName string, publicSign string, privateSign string, optKeyValue map[string]string) (*url.URL, string, error) {
	if optKeyValue == nil {
		optKeyValue = map[string]string{}
	}
	//
	//
	callbackUrlObj, _ := url.Parse(obj.config.CallbackUrl)
	callbackValue := callbackUrlObj.Query()
	callbackValue.Add("dir", dirName)
	callbackValue.Add("file", fileName)
	//
	//
	// [keys]
	//
	keys := make([]string, 0)
	for k, _ := range optKeyValue {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	propObj := miniprop.NewMiniProp()
	propObj.SetPropStringList("", "k", keys)
	//
	//
	hash := hmac.New(sha1.New, []byte(privateSign))
	//	hash := sha1.New()

	io.WriteString(hash, dirName)
	io.WriteString(hash, obj.config.Kind)
	io.WriteString(hash, fileName)
	io.WriteString(hash, publicSign)

	for _, k := range keys {
		io.WriteString(hash, optKeyValue[k])
	}
	io.WriteString(hash, string(propObj.ToJson()))
	io.WriteString(hash, privateSign)

	//
	callbackValue.Add("kv", publicSign)
	callbackValue.Add("keys", string(propObj.ToJson()))
	for k, v := range optKeyValue {
		callbackValue.Add(k, v)
	}
	callbackValue.Add("hash", base64.StdEncoding.EncodeToString(hash.Sum(nil)))
	callbackUrlObj.RawQuery = callbackValue.Encode()
	retV, retE := blobstore.UploadURL(ctx, callbackUrlObj.String(), nil)
	name := base32.StdEncoding.EncodeToString(hash.Sum([]byte("" + dirName + "/" + fileName)))
	if obj.config.HashLength >= 5 && len(name) > obj.config.HashLength {
		name = name[:obj.config.HashLength]
	}
	return retV, name, retE
}

type CheckCallbackResult struct {
	DirName  string
	FileName string
	BlobKey  string
}

func (obj *BlobManager) CheckedCallback(r *http.Request, privateSign string) (*CheckCallbackResult, error) {
	//
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		return nil, errors.New("faied to parseupload")
	}

	hashValue := r.FormValue("hash")
	dirName := r.FormValue("dir")
	fileName := r.FormValue("file")
	kv := r.FormValue("kv")

	//	hash := sha1.New()
	hash := hmac.New(sha1.New, []byte(privateSign))
	//
	// [keys]
	//
	keysSrc := r.FormValue("keys")
	propObj := miniprop.NewMiniPropFromJson([]byte(keysSrc))
	keys := propObj.GetPropStringList("", "k", make([]string, 0))
	sort.Strings(keys)
	//
	//
	io.WriteString(hash, dirName)
	io.WriteString(hash, obj.config.Kind)
	io.WriteString(hash, fileName)
	io.WriteString(hash, kv)

	for _, k := range keys {
		io.WriteString(hash, r.FormValue(k))
	}
	io.WriteString(hash, keysSrc)
	io.WriteString(hash, privateSign)

	calcHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	if 0 != strings.Compare(calcHash, hashValue) {
		return nil, errors.New("faied to check hash")
	}

	// --
	// files
	// --
	name := base32.StdEncoding.EncodeToString(hash.Sum([]byte("" + dirName + "/" + fileName)))
	if obj.config.HashLength >= 5 && len(name) > obj.config.HashLength {
		name = name[:obj.config.HashLength]
	}
	file := blobs[name]
	if len(file) == 0 {
		return nil, errors.New("faied to find file")
	}
	//
	// opt
	blobKey := string(file[0].BlobKey)
	if fileName == "" {
		fileName = blobKey
	}

	return &CheckCallbackResult{
		DirName:  dirName,
		FileName: fileName,
		BlobKey:  blobKey,
	}, nil
}
