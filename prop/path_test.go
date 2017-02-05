package prop

import (
	//	"bytes"
	"testing"
)

func TestPathA(t *testing.T) {
	pathObj := NewMiniPath("/asdf/asdf/asdf/")
	if "/asdf/asdf/asdf" != pathObj.GetDir() {
		t.Error("s")
	}

	dirs := pathObj.GetDirs()
	if "/asdf" != dirs[0] {
		t.Error("s")
	}
	if "/asdf/asdf" != dirs[1] {
		t.Error("s")
	}
	if "/asdf/asdf/asdf" != dirs[2] {
		t.Error("s")
	}
}
