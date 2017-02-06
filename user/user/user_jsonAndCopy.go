package user

import (
	"encoding/json"
	"time"

	"github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
)

// ----
// json and copy
// ----
func (userObj *User) SetUserFromsJson(ctx context.Context, source string) error {
	v := make(map[string]interface{})
	e := json.Unmarshal([]byte(source), &v)
	if e != nil {
		return e
	}
	//
	userObj.SetUserFromsMap(ctx, v)
	return nil
}

func (userObj *User) SetUserFromsMap(ctx context.Context, v map[string]interface{}) {
	propObj := prop.NewMiniPropFromMap(v)
	userObj.gaeObject.DisplayName = propObj.GetString(TypeDisplayName, "")
	userObj.gaeObject.UserName = propObj.GetString(TypeUserName, "")
	userObj.gaeObject.Created = propObj.GetTime(TypeCreated, time.Now()) //srcCreated
	userObj.gaeObject.Updated = propObj.GetTime(TypeUpdated, time.Now()) //time.Unix(0, int64(v[TypeLogined].(float64))) //srcLogin
	userObj.gaeObject.State = propObj.GetString(TypeState, "")
	userObj.gaeObject.PublicInfo = propObj.GetString(TypePublicInfo, "")
	userObj.gaeObject.PrivateInfo = propObj.GetString(TypePrivateInfo, "")
	userObj.gaeObject.Point = propObj.GetPropFloat64("", TypePoint, 0)
	userObj.gaeObject.PropValues = propObj.GetPropStringList("", TypePropValues, []string{})
	userObj.gaeObject.PropNames = propObj.GetPropStringList("", TypePropNames, []string{})
	userObj.gaeObject.IconUrl = propObj.GetString(TypeIconUrl, "")
	userObj.gaeObject.Sign = propObj.GetString(TypeSign, "")
	userObj.SetTags(propObj.GetPropStringList("", TypeTag, make([]string, 0)))
	userObj.gaeObject.Cont = propObj.GetString(TypeCont, "")
	userObj.gaeObject.Permission = propObj.GetInt(TypePermission, 0)
}

func (obj *User) ToMapPublic() map[string]interface{} {

	return map[string]interface{}{
		TypeDisplayName: obj.gaeObject.DisplayName,        //
		TypeUserName:    obj.gaeObject.UserName,           //
		TypeCreated:     obj.gaeObject.Created.UnixNano(), //
		TypeUpdated:     obj.gaeObject.Updated.UnixNano(), //
		TypeState:       obj.gaeObject.State,              //
		TypeTag:         obj.GetTags(),                    //
		TypePoint:       obj.gaeObject.Point,              //
		TypePropNames:   obj.gaeObject.PropNames,          //
		TypePropValues:  obj.gaeObject.PropValues,         //
		TypeIconUrl:     obj.gaeObject.IconUrl,            //
		TypePublicInfo:  obj.gaeObject.PublicInfo,
		TypeSign:        obj.gaeObject.Sign,
		TypeCont:        obj.gaeObject.Cont,
		TypePermission:  obj.gaeObject.Permission,
	}
}

func (obj *User) ToMapAll() map[string]interface{} {
	v := obj.ToMapPublic()
	v[TypePrivateInfo] = obj.gaeObject.PrivateInfo
	return v
}

func (obj *User) ToJson() []byte {
	return prop.NewMiniPropFromMap(obj.ToMapAll()).ToJson()
}

func (obj *User) ToJsonPublic() []byte {
	return prop.NewMiniPropFromMap(obj.ToMapPublic()).ToJson()
}
