package user

import (
	"github.com/firefirestyle/engine-v01/prop"
)

type UserKeyInfo struct {
	UserName string
	Kind     string
	Sign     string
}

func (obj *UserManager) MakeStringId(userName, sign string) string {
	propObj := prop.NewMiniProp()
	propObj.SetString("n", userName)
	propObj.SetString("k", obj.GetUserKind())
	propObj.SetString("s", sign)
	return string(propObj.ToJson())
}

func (obj *UserManager) NewUserKeyInfo(stringId string) *UserKeyInfo {
	propObj := prop.NewMiniPropFromJson([]byte(stringId))
	return &UserKeyInfo{
		UserName: propObj.GetString("n", ""),
		Kind:     propObj.GetString("k", ""),
		Sign:     propObj.GetString("s", ""),
	}
}
