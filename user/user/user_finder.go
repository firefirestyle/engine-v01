package user

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

//
//
//
func (obj *UserManager) newCursorFromSrc(cursorSrc string) *datastore.Cursor {
	c1, e := datastore.DecodeCursor(cursorSrc)
	if e != nil {
		return nil
	} else {
		return &c1
	}
}

func (obj *UserManager) makeCursorSrc(founds *datastore.Iterator) string {
	c, e := founds.Cursor()
	if e == nil {
		return c.String()
	} else {
		return ""
	}
}

//
//
func (obj *UserManager) FindUserWithUserName(ctx context.Context, userName string, keyOnly bool) *FoundUser {
	q := datastore.NewQuery(obj.config.UserKind)
	q = q.Filter("UserName =", userName) ////
	//q = q.Order("-Updated")
	q = q.Limit(obj.config.LimitOfFinding)
	return obj.FindUserFromQuery(ctx, q, "", keyOnly)
}

func (obj *UserManager) FindUserWithNewOrder(ctx context.Context, cursorSrc string, keyOnly bool) *FoundUser {
	q := datastore.NewQuery(obj.config.UserKind)
	q = q.Order("-Updated")
	q = q.Limit(obj.config.LimitOfFinding)
	return obj.FindUserFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *UserManager) FindUserWithPoint(ctx context.Context, cursorSrc string, keyOnly bool) *FoundUser {
	q := datastore.NewQuery(obj.config.UserKind)
	q = q.Order("-Point")
	q = q.Limit(obj.config.LimitOfFinding)
	return obj.FindUserFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *UserManager) FindUserFromProp(ctx context.Context, key string, value string, cursorSrc string, keyOnly bool) *FoundUser {
	q := datastore.NewQuery(obj.config.UserKind)
	v := MakePropValue(key, value)
	q = q.Filter("Props.Value =", v) ////
	q = q.Order("-Updated")
	q = q.Limit(obj.config.LimitOfFinding)
	return obj.FindUserFromQuery(ctx, q, cursorSrc, keyOnly)
}

//
//
type FoundUser struct {
	Users      []*User
	UserIds    []string
	CursorOne  string
	CursorNext string
}

func (obj *UserManager) FindUserFromQuery(ctx context.Context, queryObj *datastore.Query, cursorSrc string, keyOnly bool) *FoundUser {
	cursor := obj.newCursorFromSrc(cursorSrc)
	if cursor != nil {
		queryObj = queryObj.Start(*cursor)
	}
	queryObj = queryObj.KeysOnly()

	var userObjList []*User
	var userIdsList []string

	founds := queryObj.Run(ctx)

	var cursorNext string = ""
	var cursorOne string = ""

	for i := 0; ; i++ {
		key, err := founds.Next(nil)
		if err != nil || err == datastore.Done {
			break
		} else {
			if keyOnly == true {
				userIdsList = append(userIdsList, key.StringID())
			} else {
				userObj := obj.newUserFromStringID(ctx, key.StringID())
				errLoadUserObj := userObj.loadFromDB(ctx)
				if errLoadUserObj != nil {
					log.Infof(ctx, "Failed LoadFromDB on FindUserFromQuery "+key.StringID())
				} else {
					userObjList = append(userObjList, userObj)
					userIdsList = append(userIdsList, key.StringID())
				}
			}
		}
		if i == 0 {
			cursorOne = obj.makeCursorSrc(founds)
		}
	}
	cursorNext = obj.makeCursorSrc(founds)
	return &FoundUser{
		Users:      userObjList,
		UserIds:    userIdsList,
		CursorOne:  cursorOne,
		CursorNext: cursorNext,
	}
}
