package user

import (
	m "github.com/firefirestyle/engine-v01/prop"
)

type UserKeyInfo struct {
	UserName string
	Kind string
	Sign string
}

func (obj *UserManager) MakeUserGaeObjectKeyStringId(userName ,sign string) string {
	propObj := m.NewMiniProp()
	propObj.SetString("n", userName)
	propObj.SetString("s", sign)
	propObj.SetString("k", obj.GetUserKind())
	return string(propObj.ToJson())
}

func (obj *UserManager) GetUserKeyInfo(stringId string) *UserKeyInfo {
	propObj := m.NewMiniPropFromJson([]byte(stringId))
	return &UserKeyInfo{
		UserName: propObj.GetString("n", ""),
		Kind:propObj.GetString("k", ""),
		Sign:propObj.GetString("s", ""),
	}
}
