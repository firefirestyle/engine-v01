package user

import (
	//	"strconv"

	"time"

	"golang.org/x/net/context"
)

//
//
func (obj *UserManager) SaveUserWithImmutable(ctx context.Context, userObj *User) (*User, error) {
	// init
	//sign := strconv.Itoa(time.Now().Nanosecond())
	// copy
	userObj.SetLogined(time.Now())
	if nil != obj.SaveUser(ctx, userObj) {
		return userObj, nil
	}
	//	replayObj.SetValue(nextUserObj.GetUserName())
	return userObj, nil
}

func (obj *UserManager) GetUserFromKey(ctx context.Context, stringId string) (*User, error) {
	// Debug(ctx, "GetUserFromKey :"+stringId)
	keyInfo := obj.GetUserKeyInfo(stringId)
	// Debug(ctx, "GetUserFromKey :"+keyInfo.UserName+" : ")
	return obj.GetUserFromUserName(ctx, keyInfo.UserName)
}
