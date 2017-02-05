package blob

import (
	"time"

	"google.golang.org/appengine/datastore"
)

type GaeObjectBlobItem struct {
	RootGroup string
	Parent    string
	Name      string
	BlobKey   string
	Owner     string
	Info      string `datastore:",noindex"`
	Updated   time.Time
	Sign      string `datastore:",noindex"`
}

type BlobItem struct {
	gaeObject *GaeObjectBlobItem
	gaeKey    *datastore.Key
}

const (
	TypeRootGroup = "RootGroup"
	TypeParent    = "Parent"
	TypeName      = "Name"
	TypeBlobKey   = "BlobKey"
	TypeOwner     = "Owner"
	TypeInfo      = "Info"
	TypeUpdated   = "Updated"
	TypeSign      = "Sign"
)

func (obj *BlobItem) GetParent() string {
	return obj.gaeObject.Parent
}

func (obj *BlobItem) GetName() string {
	return obj.gaeObject.Name
}

func (obj *BlobItem) GetBlobKey() string {
	return obj.gaeObject.BlobKey
}

func (obj *BlobItem) GetInfo() string {
	return obj.gaeObject.Info
}

func (obj *BlobItem) SetInfo(v string) {
	obj.gaeObject.Info = v
}

func (obj *BlobItem) GetSign() string {
	return obj.gaeObject.Sign
}

func (obj *BlobItem) GetOwner() string {
	return obj.gaeObject.Owner
}

func (obj *BlobItem) SetOwner(v string) {
	obj.gaeObject.Owner = v
}
