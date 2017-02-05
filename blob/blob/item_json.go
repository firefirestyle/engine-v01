package blob

import (
	"encoding/json"

	"time"

	m "github.com/firefirestyle/engine-v01/prop"
)

func (obj *BlobItem) ToMap() map[string]interface{} {
	return map[string]interface{}{
		TypeRootGroup: obj.gaeObject.RootGroup,
		TypeParent:    obj.gaeObject.Parent,
		TypeName:      obj.gaeObject.Name,
		TypeBlobKey:   obj.gaeObject.BlobKey,
		TypeOwner:     obj.gaeObject.Owner,
		TypeInfo:      obj.gaeObject.Info,
		TypeUpdated:   obj.gaeObject.Updated.UnixNano(),
		TypeSign:      obj.gaeObject.Sign,
	}
}

func (obj *BlobItem) ToJson() ([]byte, error) {
	return json.Marshal(obj.ToMap())
}

func (obj *BlobItem) SetParamFromJson(source []byte) error {
	v := make(map[string]interface{})
	e := json.Unmarshal(source, &v)
	if e != nil {
		return e
	}
	//
	obj.SetParamFromMap(v)
	return nil
}

func (obj *BlobItem) SetParamFromMap(values map[string]interface{}) {
	propObj := m.NewMiniPropFromMap(values)
	obj.gaeObject.RootGroup = propObj.GetString(TypeRootGroup, "")
	obj.gaeObject.Parent = propObj.GetString(TypeParent, "")
	obj.gaeObject.Name = propObj.GetString(TypeName, "")
	obj.gaeObject.BlobKey = propObj.GetString(TypeBlobKey, "")
	obj.gaeObject.Owner = propObj.GetString(TypeOwner, "")
	obj.gaeObject.Info = propObj.GetString(TypeInfo, "")
	obj.gaeObject.Updated = propObj.GetTime(TypeUpdated, time.Now())
	obj.gaeObject.Sign = propObj.GetString(TypeSign, "")
}
