package user

import (
	"time"

	"golang.org/x/net/context"
)

func (obj *UserManager) SaveUserWithImmutable(ctx context.Context, userObj *User) (*User, error) {
	userObj.SetLogined(time.Now())
	if nil != obj.SaveUser(ctx, userObj) {
		return userObj, nil
	}
	return userObj, nil
}

func (obj *UserManager) GetUserFromKey(ctx context.Context, stringId string) (*User, error) {
	keyInfo := obj.NewUserKeyInfo(stringId)
	return obj.GetUserFromUserName(ctx, keyInfo.UserName)
}
