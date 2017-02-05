package prop

import (
	//	"bytes"
	"testing"
)

func TestPropCom(t *testing.T) {
	propObj := NewMiniProp()
	{
		propObj.SetPropInt("test", "a", 3)
		propObj.SetPropString("test", "c", "test")
		propObj.SetPropBool("a", "b", true)
		propObj.SetPropStringList("a", "c", []string{"a", "b"})
	}
	t.Log("<A>" + string(propObj.ToJson()))
	propObjB := NewMiniProp()
	{
		propObjB.SetPropInt("test", "a", 3)
		propObjB.SetPropString("test", "c", "test")
		propObjB.SetPropBool("a", "b", true)
		propObjB.SetPropStringList("a", "c", []string{"a", "b"})

	}
	propObj.CopiedOver(propObjB)
	t.Log("<B>" + string(propObj.ToJson()))
}
