package prop

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"
)

type MiniProp struct {
	prop map[string]interface{}
}

func NewMiniPropFromJsonReader(r io.Reader) *MiniProp {
	bytes, _ := ioutil.ReadAll(r)
	return NewMiniPropFromJson(bytes)
}

func NewMiniPropFromJson(source []byte) *MiniProp {
	propObj := new(MiniProp)
	propObj.prop = make(map[string]interface{})
	if source == nil {
		source = make([]byte, 0)
	}
	json.Unmarshal(source, &propObj.prop)
	return propObj
}

func NewMiniPropFromMap(source map[string]interface{}) *MiniProp {
	propObj := new(MiniProp)
	propObj.prop = source
	return propObj
}

func NewMiniProp() *MiniProp {
	propObj := new(MiniProp)
	propObj.prop = make(map[string]interface{})
	return propObj
}

func (obj *MiniProp) GetProps(category string, defaultValue interface{}) interface{} {
	if category == "" {
		return obj.prop
	} else {
		v := obj.prop[category]
		if v == nil {
			return defaultValue
		} else {
			return v
		}
	}
}

func (obj *MiniProp) Contain(key string) bool {
	return obj.ContainProp("", key)
}

func (obj *MiniProp) ContainProp(category string, key string) bool {
	if category == "" {
		_, ok := obj.prop[key]
		return ok
	} else {
		v := obj.prop[category]
		if v == nil {
			return false
		}
		if obj.prop[category].(map[string]interface{})[key] == nil {
			return false
		}
		_, ok := obj.prop[category].(map[string]interface{})[key]
		return ok
	}
}

func (obj *MiniProp) GetProp(category string, key string, defaultValue interface{}) interface{} {
	if category == "" {
		v := obj.prop[key]
		if v == nil {
			return defaultValue
		} else {
			return v
		}
	} else {
		v := obj.prop[category]
		if v == nil {
			return defaultValue
		}
		if obj.prop[category].(map[string]interface{})[key] == nil {
			return defaultValue
		}
		return obj.prop[category].(map[string]interface{})[key]
	}
}

func (obj *MiniProp) GetPropBool(category string, key string, defaultValue bool) bool {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case bool:
		return obj.GetProp(category, key, defaultValue).(bool)
	}
	return defaultValue
}

func (obj *MiniProp) GetPropInt(category string, key string, defaultValue int) int {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case int:
		return obj.GetProp(category, key, defaultValue).(int)
	case float64:
		return int(obj.GetProp(category, key, defaultValue).(float64))
	}
	return defaultValue
}

func (obj *MiniProp) GetPropFloat64(category string, key string, defaultValue float64) float64 {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case float64:
		return obj.GetProp(category, key, defaultValue).(float64)
	}
	return defaultValue
}

func (obj *MiniProp) GetPropString(category string, key string, defaultValue string) string {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case string:
		return obj.GetProp(category, key, defaultValue).(string)
	}
	return defaultValue
}

func (obj *MiniProp) GetPropBytes(category string, key string, defaultValue []byte) []byte {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case string:
		va, er := base64.StdEncoding.DecodeString(v.(string))
		if er == nil {
			return va
		} else {
			return defaultValue
		}
	}
	return defaultValue
}

func (obj *MiniProp) GetPropTime(category string, key string, defaultValue time.Time) time.Time {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case float64:
		return time.Unix(0, int64(v.(float64)))
	}
	return defaultValue
}

func (obj *MiniProp) SetProp(category string, key string, value interface{}) {
	if category == "" {
		obj.prop[key] = value
	} else {
		v := obj.prop[category]
		if v == nil {
			obj.prop[category] = make(map[string]interface{})
		}
		obj.prop[category].(map[string]interface{})[key] = value
	}
}

func (obj *MiniProp) SetPropBool(category string, key string, value bool) {
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) SetPropInt(category string, key string, value int) {
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) SetPropFloat(category string, key string, value float64) {
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) SetPropString(category string, key string, value string) {
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) SetPropBytes(category string, key string, value []byte) {
	obj.SetProp(category, key, base64.StdEncoding.EncodeToString([]byte(value)))
}

func (obj *MiniProp) SetPropTime(category string, key string, value time.Time) {
	obj.SetProp(category, key, value.UnixNano())
}

func (obj *MiniProp) ToJson() []byte {
	v, e := json.Marshal(obj.prop)
	if e != nil {
		return []byte("{}")
	}
	return v
}

func (obj *MiniProp) ToMap() map[string]interface{} {
	return obj.prop
}

func (obj *MiniProp) ToJsonFromCategory(category string) []byte {
	v := obj.GetProps(category, make(map[string]interface{}))
	vv, e := json.Marshal(v)
	if e != nil {
		return []byte("{}")
	}
	return vv
}

// ---
// Single
// ---
func (obj *MiniProp) GetString(key string, defaultValue string) string {
	return obj.GetPropString("", key, defaultValue)
}

func (obj *MiniProp) SetString(key string, value string) {
	obj.SetPropString("", key, value)
}

func (obj *MiniProp) GetBool(key string, defaultValue bool) bool {
	return obj.GetPropBool("", key, defaultValue)
}

func (obj *MiniProp) SetBool(key string, value bool) {
	obj.SetPropBool("", key, value)
}

func (obj *MiniProp) GetInt(key string, defaultValue int) int {
	return obj.GetPropInt("", key, defaultValue)
}

func (obj *MiniProp) SetInt(key string, value int) {
	obj.SetPropInt("", key, value)
}

func (obj *MiniProp) GetFloat(key string, defaultValue float64) float64 {
	return obj.GetPropFloat64("", key, defaultValue)
}

func (obj *MiniProp) SetFloatString(key string, value float64) {
	obj.SetPropFloat("", key, value)
}

func (obj *MiniProp) GetTime(key string, defaultValue time.Time) time.Time {
	return obj.GetPropTime("", key, defaultValue)
}

func (obj *MiniProp) SetTime(key string, value time.Time) {
	obj.SetPropTime("", key, value)
}

// ---
// List
// ---
func (obj *MiniProp) SetPropStringList(category string, key string, value []string) {
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) SetPropIntList(category string, key string, value []int) {
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) SetPropFloatList(category string, key string, value []float64) {
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) SetPropBoolList(category string, key string, value []bool) {
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) GetPropStringList(category string, key string, defaultValue []string) []string {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case []string:
		return v.([]string)
	case []interface{}:
		vs := v.([]interface{})
		ret := make([]string, 0)
		for _, vv := range vs {
			switch vv.(type) {
			case string:
				ret = append(ret, vv.(string))
			default:
				return defaultValue
			}
		}
		return ret
	}
	return defaultValue
}

func (obj *MiniProp) GetPropFloatList(category string, key string, defaultValue []float64) []float64 {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case []float64:
		return v.([]float64)
	case []interface{}:
		vs := v.([]interface{})
		ret := make([]float64, 0)
		for _, vv := range vs {
			switch vv.(type) {
			case float64:
				ret = append(ret, vv.(float64))
			default:
				return defaultValue
			}
		}
		return ret
	}
	return defaultValue
}

func (obj *MiniProp) GetPropBoolList(category string, key string, defaultValue []bool) []bool {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case []bool:
		return v.([]bool)
	case []interface{}:
		vs := v.([]interface{})
		ret := make([]bool, 0)
		for _, vv := range vs {
			switch vv.(type) {
			case bool:
				ret = append(ret, vv.(bool))
			default:
				return defaultValue
			}
		}
		return ret
	}
	return defaultValue
}

func (obj *MiniProp) GetPropIntList(category string, key string, defaultValue []int) []int {
	v := obj.GetProp(category, key, defaultValue)
	switch v.(type) {
	case []int:
		return v.([]int)
	case []interface{}:
		{
			vs := v.([]interface{})
			ret := make([]int, 0)
			for _, vv := range vs {
				switch vv.(type) {
				case int:
					ret = append(ret, vv.(int))
				case float64:
					ret = append(ret, int(vv.(float64)))
				default:
					return defaultValue
				}
			}
			return ret
		}
	}
	return defaultValue
}

// ---
// Double
// ---
func (obj *MiniProp) AddPropListItem(category string, key string, value interface{}) {
	v := obj.GetProp(category, key, nil)
	if v == nil {
		v = make([]interface{}, 0)
	}

	switch v.(type) {
	case []interface{}:
	default:
		v = make([]interface{}, 0)
	}
	v = append(v.([]interface{}), value)
	obj.SetProp(category, key, value)
}

func (obj *MiniProp) GetLengthPropListItem(category string, key string) int {
	v := obj.GetProp(category, key, nil)
	if v == nil {
		return 0
	}

	switch v.(type) {
	case []interface{}:
		return len(v.([]interface{}))
	default:
		return 0
	}
}

func (obj *MiniProp) GetPropListItem(category string, key string, index int, defaultValue interface{}) interface{} {
	v := obj.GetProp(category, key, nil)
	if v == nil {
		return defaultValue
	}

	switch v.(type) {
	case []interface{}:
		if index < len(v.([]interface{})) {
			return v.([]interface{})[index]
		} else {
			return defaultValue
		}
	default:
		return defaultValue
	}
}
