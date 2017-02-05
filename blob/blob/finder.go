package blob

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

type BlobFounds struct {
	Keys       []string
	CursorNext string
	CursorOne  string
}

/*
https://cloud.google.com/appengine/docs/go/config/indexconfig#updating_indexes
*/
func (obj *BlobManager) FindBlobItemFromParent(ctx context.Context, parent string, cursorSrc string) BlobFounds {
	//
	q := datastore.NewQuery(obj.config.Kind)
	q = q.Filter("RootGroup =", obj.config.RootGroup)
	q = q.Filter("Parent =", parent)
	q = q.Order("-Updated")
	//
	return obj.FindBlobItemFromQuery(ctx, q, cursorSrc)
}

func (obj *BlobManager) FindBlobItemFromPath(ctx context.Context, parent string, name string, cursorSrc string) BlobFounds {
	//
	q := datastore.NewQuery(obj.config.Kind)
	q = q.Filter("RootGroup =", obj.config.RootGroup)
	q = q.Filter("Parent =", parent)
	q = q.Filter("Name =", name)
	q = q.Order("-Updated")
	//
	return obj.FindBlobItemFromQuery(ctx, q, cursorSrc)
}

func (obj *BlobManager) FindAllBlobItemFromPath(ctx context.Context, parent string) BlobFounds {
	//
	q := datastore.NewQuery(obj.config.Kind)
	q = q.Filter("RootGroup =", obj.config.RootGroup)
	q = q.Filter("Parent =", parent)
	q = q.Order("-Updated")
	//
	return obj.FindBlobItemFromQueryAll(ctx, q)
}

func (obj *BlobManager) FindBlobItemFromOwner(ctx context.Context, owner string, cursorSrc string) BlobFounds {
	//
	q := datastore.NewQuery(obj.config.Kind)
	q = q.Filter("RootGroup =", obj.config.RootGroup)
	q = q.Filter("Owner =", owner)
	q = q.Order("-Updated")
	//
	return obj.FindBlobItemFromQuery(ctx, q, cursorSrc)
}

//
//
func (obj *BlobManager) FindBlobItemFromQueryAll(ctx context.Context, q *datastore.Query) BlobFounds {
	founded := obj.FindBlobItemFromQuery(ctx, q, "")
	oneCursor := founded.CursorOne
	nextCursor := founded.CursorNext
	keys := make([]string, 0)
	for {
		if len(founded.Keys) <= 0 {
			break
		}
		for _, v := range founded.Keys {
			keys = append(keys, v)
		}
		prevFounded := founded
		founded = obj.FindBlobItemFromQuery(ctx, q, nextCursor)
		nextCursor = founded.CursorNext
		if prevFounded.CursorOne == founded.CursorOne {
			break
		}
	}
	return BlobFounds{
		Keys:       keys,
		CursorNext: nextCursor,
		CursorOne:  oneCursor,
	}
}

func (obj *BlobManager) FindBlobItemFromQuery(ctx context.Context, q *datastore.Query, cursorSrc string) BlobFounds {
	cursor := obj.newCursorFromSrc(cursorSrc)
	if cursor != nil {
		q = q.Start(*cursor)
	}
	q = q.KeysOnly()
	founds := q.Run(ctx)

	var keys []string
	var cursorNext string = ""
	var cursorOne string = ""

	for i := 0; ; i++ {
		key, err := founds.Next(nil)
		if err != nil || err == datastore.Done {
			break
		} else {
			keys = append(keys, key.StringID())
		}
		if i == 0 {
			cursorOne = obj.makeCursorSrc(founds)
		}
	}
	cursorNext = obj.makeCursorSrc(founds)
	return BlobFounds{
		Keys:       keys,
		CursorOne:  cursorOne,
		CursorNext: cursorNext,
	}
}

func (obj *BlobManager) newCursorFromSrc(cursorSrc string) *datastore.Cursor {
	c1, e := datastore.DecodeCursor(cursorSrc)
	if e != nil {
		return nil
	} else {
		return &c1
	}
}

func (obj *BlobManager) makeCursorSrc(founds *datastore.Iterator) string {
	c, e := founds.Cursor()
	if e == nil {
		return c.String()
	} else {
		return ""
	}
}
