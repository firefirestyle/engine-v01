package prop

import (
	"strings"
)

type MiniPath struct {
	dir string
}

func NewMiniPath(dir string) *MiniPath {
	return &MiniPath{
		dir: dir,
	}
}

func (obj *MiniPath) GetDir() string {
	items := strings.Split(obj.dir, "/")
	buf := make([]byte, 0, len(items)*2)
	for _, v := range items {
		if v != "" {
			buf = append(buf, "/"...)
			buf = append(buf, v...)
		}
	}
	return string(buf)
}

func (obj *MiniPath) GetDirs() []string {
	items := strings.Split(obj.dir, "/")
	ret := make([]string, 0, len(items)*2)
	buf := make([]byte, 0, len(items)*2)
	for _, v := range items {
		if v != "" {
			buf = append(buf, "/"...)
			buf = append(buf, v...)
			ret = append(ret, string(buf))
		}
	}
	return ret
}
