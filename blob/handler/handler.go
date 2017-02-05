package handler

import (
	"net/http"

	mm "github.com/firefirestyle/engine-v01/prop"

	miniblob "github.com/firefirestyle/engine-v01/blob/blob"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

type BlobHandlerOnEvent struct {
	OnBlobRequestList    []func(w http.ResponseWriter, r *http.Request, input *mm.MiniProp, output *mm.MiniProp, h *BlobHandler) (map[string]string, error)
	OnBlobBeforeSaveList []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem) error
	OnBlobCompleteList   []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem) error
	OnBlobFailedList     []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem)
	OnDeleteRequestList  []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler) error
	OnDeleteFailedList   []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem)
	OnDeleteSuccessList  []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem)
	OnGetRequestList     []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler) error
	OnGetFailedList      []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem)
	OnGetSuccessList     []func(http.ResponseWriter, *http.Request, *mm.MiniProp, *BlobHandler, *miniblob.BlobItem)
}

type BlobHandler struct {
	manager     *miniblob.BlobManager
	onEvent     BlobHandlerOnEvent
	callbackUrl string
	privateSign string
}

//
//
//
//func (obj *BlobHandler) GetBlobHandleEvent() *BlobHandlerOnEvent {
//	return &obj.onEvent
//}
//
//
//
//
func (obj *BlobHandler) GetManager() *miniblob.BlobManager {
	return obj.manager
}

func NewBlobHandler(callbackUrl string, privateSign string, config miniblob.BlobManagerConfig) *BlobHandler {
	handlerObj := new(BlobHandler)
	handlerObj.privateSign = privateSign
	handlerObj.callbackUrl = callbackUrl
	handlerObj.manager = miniblob.NewBlobManager(config)
	handlerObj.onEvent = BlobHandlerOnEvent{}
	return handlerObj
}

func HandleError(w http.ResponseWriter, r *http.Request, outputProp *mm.MiniProp, errorCode int, errorMessage string) {
	//
	//
	if errorCode != 0 {
		outputProp.SetInt("errorCode", errorCode)
	}
	if errorMessage != "" {
		outputProp.SetString("errorMessage", errorMessage)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(outputProp.ToJson())
}

const (
	ErrorCodeAtGetRequestCheck             = 2001
	ErrorCodeAtGetRequestFindBlobItem      = 2002
	ErrorCodeAtDeleteRequestCheck          = 2001
	ErrorCodeAtDeleteRequestFindBlobItem   = 2002
	ErrorCodeAtDeleteRequestDeleteBlobItem = 2003
	ErrorCodeAtBlobRequestCheck            = 2001
	ErrorCodeAtBlobMakeRequestUrl          = 2004
	ErrorCodeAtBlobCheckCallback           = 2005
	ErrorCodeAtBlobBeforeSaveCheck         = 2007
	ErrorCodeAtBlobCompleteCheck           = 2008
	ErrorCodeAtBlobSaveBlobItem            = 2009
	ErrorCodeGetBlobItem                   = 2010
	ErrorCodeDeleteBlobItem                = 2011
)

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
