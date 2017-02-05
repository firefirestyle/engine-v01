package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	miniprop "github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func (obj *ArticleHandler) SaveData(ctx context.Context, token string, path string, articleId string, dirName string, fileName string, data []byte) error {
	//
	url, name, err := obj.GetRequestCodeToSaveData(ctx, token, path, articleId, dirName, fileName)
	if err != nil {
		return err
	}
	//
	//
	Debug(ctx, "--2--"+url)
	var b bytes.Buffer
	fw := multipart.NewWriter(&b)
	file, err := fw.CreateFormFile(name, "blob")
	if err != nil {
		Debug(ctx, "--2--")
		return err
	}
	if _, err = file.Write(data); err != nil {
		Debug(ctx, "--3--")
		return err
	}
	fw.Close()
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		Debug(ctx, "--4--")
		return err
	}
	//propName, "blob"
	req.Header.Set("Content-Type", fw.FormDataContentType())

	//
	//
	client := urlfetch.Client(ctx)
	res, err := client.Do(req)
	if err != nil {
		Debug(ctx, "--5--")
		return err
	}

	if res.StatusCode != 200 {
		Debug(ctx, "--6-- :"+res.Status)
		return err
	}
	return nil
}

func (obj *ArticleHandler) GetRequestCodeToSaveData(ctx context.Context, token string, path string, articleId string, dirName string, fileName string) (string, string, error) {
	//
	v := appengine.DefaultVersionHostname(ctx)
	scheme := "https"
	if v == "127.0.0.1:8080" || v == "localhost:8080" {
		v = "localhost:8080"
		scheme = "http"
	}
	baseUrl := fmt.Sprintf("%s://%s%s", scheme, v, path)
	//
	client := urlfetch.Client(ctx)
	//	token

	propObj := miniprop.NewMiniProp()
	propObj.SetString("token", token)
	propObj.SetString("articleId", articleId)
	propObj.SetString("dir", dirName)
	propObj.SetString("file", fileName)
	res, err := client.Post(baseUrl, "application/json", bytes.NewReader(propObj.ToJson()))
	if err != nil {
		return "", "", err
	}
	vv, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return "", "", err2
	}
	outPropObj := miniprop.NewMiniPropFromJson(vv)
	return outPropObj.GetString("token", ""), outPropObj.GetString("name", ""), nil
}

/*
func (obj *BlobManager) saveData(ctx context.Context, url string, data []byte) error {
	Debug(ctx, "--2--"+url)
	var b bytes.Buffer
	fw := multipart.NewWriter(&b)
	file, err := fw.CreateFormField("file")
	if err != nil {
		Debug(ctx, "--2--")
		return err
	}
	if _, err = file.Write(data); err != nil {
		Debug(ctx, "--3--")
		return err
	}
	fw.Close()
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		Debug(ctx, "--4--")
		return err
	}
	req.Header.Set("Content-Type", fw.FormDataContentType())

	//
	//
	client := urlfetch.Client(ctx)
	res, err := client.Do(req)
	if err != nil {
		Debug(ctx, "--5--")
		return err
	}

	//http.NewRequest("GET")
	//blobs, _, err := blobstore.ParseUpload(r)
	if res.StatusCode != http.StatusCreated {
		//if err != nil {
		v, _ := ioutil.ReadAll(res.Body)
		Debug(ctx, "--6--"+fmt.Sprintf("%d %s %s", res.StatusCode, res.Status, string(v)))
		return err
	}
	//	blobs.
	//	Debug(ctx, "--6--"+fmt.Sprintf("%d %s", blobs.))

	return nil
}
*/
