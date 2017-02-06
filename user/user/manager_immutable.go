package user

import (
	"golang.org/x/net/context"
)

func (obj *UserManager) SaveUserWithImmutable(ctx context.Context, userObj *User) (*User, error) {
	return obj.SaveUser(ctx, userObj)
}

func (obj *UserManager) GetUserFromKey(ctx context.Context, stringId string) (*User, error) {
	keyInfo := obj.NewUserKeyInfo(stringId)
	return obj.GetUserFromUserName(ctx, keyInfo.UserName)
}
