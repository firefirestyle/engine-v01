package article

import (
	miniprop "github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

//
//
//
type ArtQuery struct {
	q *datastore.Query
}

func (obj *ArticleManager) NewArtQuery() *ArtQuery {
	return &ArtQuery{
		q: datastore.NewQuery(obj.config.KindArticle),
	}

}

func (obj *ArtQuery) GetQuery() *datastore.Query {
	return obj.q
}

func (obj *ArtQuery) WithProp(ctx context.Context, props map[string]string) *ArtQuery {
	for k, v := range props {
		p := miniprop.NewMiniProp()
		p.SetString(k, v)
		v := string(p.ToJson())
		obj.q = obj.q.Filter("Props.Value =", v) ////
	}
	return obj
}

func (obj *ArtQuery) WithTags(ctx context.Context, tags []string) *ArtQuery {
	for _, tag := range tags {
		obj.q = obj.q.Filter("Tags.Tag =", tag) ////
	}
	return obj
}

func (obj *ArtQuery) WithUserName(ctx context.Context, userName string) *ArtQuery {
	obj.q = obj.q.Filter("UserName =", userName)
	return obj
}

func (obj *ArtQuery) WithUpdateMinus(ctx context.Context) *ArtQuery {
	obj.q = obj.q.Order("-Updated")
	return obj
}
func (obj *ArtQuery) WithUpdatePulas(ctx context.Context) *ArtQuery {
	obj.q = obj.q.Order("Updated")
	return obj
}

func (obj *ArtQuery) WithLimitOfFinding(ctx context.Context, limitOfFinding int) *ArtQuery {
	obj.q = obj.q.Limit(limitOfFinding)
	return obj
}
