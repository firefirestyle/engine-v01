package user

import (
	"errors"

	"github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type UserManagerConfig struct {
	UserKind        string
	UserPointerKind string
	LengthHash      int
	LimitOfFinding  int
}

type UserManager struct {
	config UserManagerConfig
}

func NewUserManager(config UserManagerConfig) *UserManager {
	obj := new(UserManager)
	if config.UserKind == "" {
		config.UserKind = "fu"
	}
	if config.UserPointerKind == "" {
		config.UserPointerKind = config.UserKind + "-pointer"
	}
	if config.LimitOfFinding <= 0 {
		config.LimitOfFinding = 20
	}
	obj.config = config

	return obj
}

func (obj *UserManager) GetUserKind() string {
	return obj.config.UserKind
}

func (obj *UserManager) NewNewUser(ctx context.Context) *User {
	return obj.newUserWithUserName(ctx)
}

func (obj *UserManager) GetUserFromUserName(ctx context.Context, userName string) (*User, error) {
	foundUser := obj.FindUserWithUserName(ctx, userName, false)
	if len(foundUser.Users) == 0 {
		return nil, errors.New("Not found " + userName)
	}
	return foundUser.Users[0], nil
}

func (obj *UserManager) SaveUser(ctx context.Context, userObj *User) (*User, error) {
	nextUser := obj.cloneUser(ctx, userObj)
	e := nextUser.pushToDB(ctx)
	if e == nil {
		datastore.Delete(ctx, userObj.gaeObjectKey)
		return nextUser, e
	} else {
		return nil, e
	}
}

func (obj *UserManager) DeleteUser(ctx context.Context, userName string, sign string) error {
	gaeKey := obj.newUserGaeObjectKey(ctx, userName, sign)
	return datastore.Delete(ctx, gaeKey)
}

//
func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}

func MakePropValue(name, v string) string {
	p := prop.NewMiniProp()
	p.SetString(name, v)
	v = string(p.ToJson())
	return v
}
