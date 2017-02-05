package prop

import (
	"bytes"
	"testing"
)

func TestKey(t *testing.T) {
	propObj := NewMiniProp()
	propObj.SetPropInt("test", "a", 3)
	propObj.SetPropString("test", "c", "test")
	if 1 != propObj.GetPropInt("d", "b", 1) {
		t.Error("s")
	}
	//
	//
	propObj = NewMiniPropFromJson(propObj.ToJson())
	if 1 != propObj.GetPropInt("d", "b", 1) {
		t.Error("s")
	}
}

func TestOne(t *testing.T) {
	propObj := NewMiniProp()
	propObj.SetBool("a", true)
	propObj.SetInt("b", 1)
	propObj.SetString("c", "test")

	if true != propObj.GetBool("a", false) {
		t.Error("s")
	}
	if 1 != propObj.GetInt("b", 0) {
		t.Error("s")
	}
	if "test" != propObj.GetString("c", "") {
		t.Error("s")
	}
	if "{\"a\":true,\"b\":1,\"c\":\"test\"}" != string(propObj.ToJson()) {
		t.Error("s")
	}
}

func TestBool(t *testing.T) {
	propObj := NewMiniProp()
	propObj.SetPropBool("test", "a", false)
	propObj.SetPropString("test", "c", "test")
	if false != propObj.GetPropBool("test", "a", true) {
		t.Error("s")
	}
	if 1 != propObj.GetPropInt("test", "b", 1) {
		t.Error("s")
	}
	if 1 != propObj.GetPropInt("test", "c", 1) {
		t.Error("s")
	}
	//
	//
	propObj = NewMiniPropFromJson(propObj.ToJson())
	if false != propObj.GetPropBool("test", "a", true) {
		t.Error("s")
	}
	if 1 != propObj.GetPropInt("test", "b", 1) {
		t.Error("s")
	}
	if 1 != propObj.GetPropInt("test", "c", 1) {
		t.Error("s")
	}
}

func TestInt(t *testing.T) {
	propObj := NewMiniProp()
	propObj.SetPropInt("test", "a", 3)
	propObj.SetPropString("test", "c", "test")
	if 3 != propObj.GetPropInt("test", "a", 0) {
		t.Error("s")
	}
	if 1 != propObj.GetPropInt("test", "b", 1) {
		t.Error("s")
	}
	if 1 != propObj.GetPropInt("test", "c", 1) {
		t.Error("s")
	}
	//
	//
	propObj = NewMiniPropFromJson(propObj.ToJson())
	if 3 != propObj.GetPropInt("test", "a", 0) {
		t.Error("s")
	}
	if 1 != propObj.GetPropInt("test", "b", 1) {
		t.Error("s")
	}
	if 1 != propObj.GetPropInt("test", "c", 1) {
		t.Error("s")
	}
}

func TestFloat(t *testing.T) {
	propObj := NewMiniProp()
	propObj.SetPropFloat("test", "a", 3.0)
	propObj.SetPropString("test", "c", "test")
	if 3.0 != propObj.GetPropFloat64("test", "a", 0) {
		t.Error("s1")
	}
	if 1 != propObj.GetPropInt("test", "b", 1) {
		t.Error("s2")
	}
	if 1 != propObj.GetPropInt("test", "c", 1) {
		t.Error("s3")
	}
	//
	//
	propObj = NewMiniPropFromJson(propObj.ToJson())
	if 3.0 != propObj.GetPropFloat64("test", "a", 0) {
		t.Error("s1")
	}
	if 1 != propObj.GetPropInt("test", "b", 1) {
		t.Error("s5")
	}
	if 1 != propObj.GetPropInt("test", "c", 1) {
		t.Error("s6z")
	}
}

func TestString(t *testing.T) {
	propObj := NewMiniProp()
	propObj.SetPropString("test", "a", "3")
	propObj.SetPropInt("test", "c", 1)
	if "3" != propObj.GetPropString("test", "a", "0") {
		t.Error("s")
	}
	if "1" != propObj.GetPropString("test", "b", "1") {
		t.Error("s")
	}
	if "1" != propObj.GetPropString("test", "c", "1") {
		t.Error("s")
	}
	//
	//
	propObj = NewMiniPropFromJson(propObj.ToJson())
	if "3" != propObj.GetPropString("test", "a", "0") {
		t.Error("s")
	}
	if "1" != propObj.GetPropString("test", "b", "1") {
		t.Error("s")
	}
	if "1" != propObj.GetPropString("test", "c", "1") {
		t.Error("s")
	}
}

func TestBytes(t *testing.T) {
	propObj := NewMiniProp()
	propObj.SetPropBytes("test", "a", []byte{1, 2, 3})
	propObj.SetPropInt("test", "c", 1)

	if 0 != bytes.Compare([]byte{1, 2, 3}, propObj.GetPropBytes("test", "a", []byte{1, 2})) {
		t.Error("TestBytes1")
	}
	if 0 != bytes.Compare([]byte{1, 2}, propObj.GetPropBytes("test", "b", []byte{1, 2})) {
		t.Error("TestBytes2")
	}
	if 0 != bytes.Compare([]byte{1, 2}, propObj.GetPropBytes("test", "c", []byte{1, 2})) {
		t.Error("TestBytes3")
	}
	//
	//
	propObj = NewMiniPropFromJson(propObj.ToJson())
	if 0 != bytes.Compare([]byte{1, 2, 3}, propObj.GetPropBytes("test", "a", []byte{1, 2})) {
		t.Error("TestBytes1")
	}
	if 0 != bytes.Compare([]byte{1, 2}, propObj.GetPropBytes("test", "b", []byte{1, 2})) {
		t.Error("TestBytes2")
	}
	if 0 != bytes.Compare([]byte{1, 2}, propObj.GetPropBytes("test", "c", []byte{1, 2})) {
		t.Error("TestBytes3")
	}
}

func TestStringList(t *testing.T) {
	propObj := NewMiniProp()
	propObj.SetPropStringList("test", "a", []string{"1", "2", "3"})
	ret := propObj.GetPropStringList("test", "a", []string{"1", "1", "1"})
	if !(ret[0] == "1" && ret[1] == "2" && ret[2] == "3") {
		t.Error("TestBytes3")
	}
	propObj = NewMiniPropFromJson(propObj.ToJson())
	ret = propObj.GetPropStringList("test", "a", []string{"1", "1", "1"})
	//	t.Log(">>>>>" + len(ret1))
	if !(ret[0] == "1" && ret[1] == "2" && ret[2] == "3") {
		t.Error("TestBytes4" + string(propObj.ToJson()) + ":")
	}

}
