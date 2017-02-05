package sns

import (
	"net/http"
	"net/url"

	"errors"

	"crypto/sha1"
	"encoding/base64"
	"io"

	"strings"

	m "github.com/firefirestyle/engine-v01/prop"

	"sort"

	"crypto/hmac"
	"strconv"
	"time"
)

//
//
//
func WithHashAndValue(tmpValues url.Values, privateSign string, callbackUrl string, opts map[string]string) {
	publicSign := strconv.FormatInt(time.Now().Unix(), 36)
	keys := make([]string, 0)
	{
		for k, _ := range opts {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}
	propObj := m.NewMiniProp()
	propObj.AddPropListItem("", "ks", keys)
	hash := hmac.New(sha1.New, []byte(privateSign))
	io.WriteString(hash, publicSign)
	io.WriteString(hash, string(propObj.ToJson()))
	io.WriteString(hash, callbackUrl)
	//
	for _, k := range keys {
		io.WriteString(hash, opts[k])
	}
	//
	io.WriteString(hash, privateSign)
	calcHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	tmpValues.Add("ps", publicSign)
	tmpValues.Add("hash", calcHash)
	tmpValues.Add("ks", string(propObj.ToJson()))
}

func CheckHashAndValue(w http.ResponseWriter, r *http.Request, secretSign string) error {
	//
	values := r.URL.Query()
	hashV := values.Get("hash")
	publicSign := values.Get("ps")
	ks := values.Get("ks")
	clCallback := values.Get("cb")
	propObj := m.NewMiniPropFromJson([]byte(ks))
	keys := propObj.GetPropStringList("", "ks", nil)
	{
		hash := hmac.New(sha1.New, []byte(secretSign))
		io.WriteString(hash, publicSign)
		io.WriteString(hash, ks)
		io.WriteString(hash, clCallback)
		for _, v := range keys {
			io.WriteString(hash, r.FormValue(v))
		}
		io.WriteString(hash, secretSign)
		calcHash := base64.StdEncoding.EncodeToString(hash.Sum(nil))

		if strings.Compare(calcHash, hashV) != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return errors.New("Wrong Hash")
		}
	}
	return nil
}
