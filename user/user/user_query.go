package user

import (
	miniprop "github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

//
//
//
type UserQuery struct {
	q *datastore.Query
}

func (obj *UserManager) NewUserQuery() *UserQuery {
	return &UserQuery{
		q: datastore.NewQuery(obj.config.UserKind),
	}
}

func (obj *UserQuery) GetQuery() *datastore.Query {
	return obj.q
}

func (obj *UserQuery) WithProp(ctx context.Context, props map[string]string) *UserQuery {
	for k, v := range props {
		p := miniprop.NewMiniProp()
		p.SetString(k, v)
		v := string(p.ToJson())
		obj.q = obj.q.Filter("Props.Value =", v) ////
	}
	return obj
}

func (obj *UserQuery) WithTags(ctx context.Context, tags []string) *UserQuery {
	for _, tag := range tags {
		obj.q = obj.q.Filter("Tags.Tag =", tag) ////
	}
	return obj
}

func (obj *UserQuery) WithUserName(ctx context.Context, userName string) *UserQuery {
	obj.q = obj.q.Filter("UserName =", userName)
	return obj
}

func (obj *UserQuery) WithStatus(ctx context.Context, status string) *UserQuery {
	obj.q = obj.q.Filter("State =", status)
	return obj
}

func (obj *UserQuery) WithUpdateMinus(ctx context.Context) *UserQuery {
	obj.q = obj.q.Order("-Updated")
	return obj
}

func (obj *UserQuery) WithUpdatePulas(ctx context.Context) *UserQuery {
	obj.q = obj.q.Order("Updated")
	return obj
}

func (obj *UserQuery) WithLimitOfFinding(ctx context.Context, limitOfFinding int) *UserQuery {
	obj.q = obj.q.Limit(limitOfFinding)
	return obj
}
