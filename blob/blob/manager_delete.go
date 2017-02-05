package blob

import (
	m "github.com/firefirestyle/engine-v01/prop"
	"golang.org/x/net/context"
)

func (obj *BlobManager) SaveBlobItemWithImmutable(ctx context.Context, newItem *BlobItem) error {
	//
	// mkdirs
	pathObj := m.NewMiniPath(newItem.GetParent())
	_, parentDirErr := obj.GetBlobItem(ctx, pathObj.GetDir(), ".dir", "")
	if parentDirErr != nil {
		for _, v := range pathObj.GetDirs() {
			dirObj := obj.NewBlobItem(ctx, v, ".dir", "")
			dirErr := dirObj.saveDB(ctx)
			if dirErr != nil {
				return dirErr
			}
		}
	}
	//
	//
	blobStringId, _, currErr := obj.GetBlobItemStringIdFromPointer(ctx, newItem.GetParent(), newItem.GetName())

	errSave := newItem.saveDB(ctx)
	if errSave != nil {
		return errSave
	}

	if currErr == nil {
		err := obj.DeleteBlobItemFromStringId(ctx, blobStringId)
		if err != nil {
			Debug(ctx, "<gomidata>"+blobStringId+"</gomidata>")
		}
	}
	obj.SaveSignCache(ctx, newItem.GetParent(), newItem.GetName(), newItem.GetSign())
	return nil

}

func (obj *BlobManager) DeleteBlobItem(ctx context.Context, item *BlobItem) error {
	obj.SaveSignCache(ctx, item.GetParent(), item.GetName(), "")
	return obj.DeleteBlobItemFromStringId(ctx, item.gaeKey.StringID())
}
func (obj *BlobManager) DeleteBlobItemWithPointerFromStringId(ctx context.Context, blolStringId string) error {
	idInfo := obj.GetKeyInfoFromStringId(blolStringId)
	obj.SaveSignCache(ctx, idInfo.Parent, idInfo.Name, "")
	return obj.DeleteBlobItemFromStringId(ctx, blolStringId)
}

//
//
func (obj *BlobManager) DeleteBlobItemsWithPointerAtRecursiveMode(ctx context.Context, parent string) error {
	folders := make([]string, 0)
	folders = append(folders, parent)
	foldersTmp := make([]string, 0)
	for len(folders) > 0 {
		folder := folders[0]
		folders = folders[1:]
		foldersTmp = append(foldersTmp, folder)
		//
		founded := obj.FindAllBlobItemFromPath(ctx, folder)
		for _, v := range founded.Keys {
			keyInfo := obj.GetKeyInfoFromStringId(v)
			if keyInfo.Name == ".dir" {
				folders = append(folders, v)
				continue
			}
			blobObj, blobErr := obj.GetBlobItem(ctx, keyInfo.Parent, keyInfo.Name, keyInfo.Sign)
			if blobErr == nil {
				obj.DeleteBlobItem(ctx, blobObj)
			}
		}
	}
	for _, v := range foldersTmp {
		keyInfo := obj.GetKeyInfoFromStringId(v)
		blobObj, blobErr := obj.GetBlobItem(ctx, keyInfo.Parent, keyInfo.Name, keyInfo.Sign)
		if blobErr == nil {
			obj.DeleteBlobItem(ctx, blobObj)
		}
	}
	return nil
}
