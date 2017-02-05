package prop

//import (
//	"reflect"
//)

func (obj *MiniProp) CopiedOver(b *MiniProp) {
	//	reflect.Copy()copy(obj.prop, b.prop)
	obj.CopiedOverMaps(obj.prop, b.prop)
}

func (obj *MiniProp) CopiedOverMaps(propA, propB map[string]interface{}) {
	for k, v := range propB {
		if propA[k] == nil {
			propA[k] = v
			continue
		}
		switch propA[k].(type) {
		case bool, int, int8, int16, int32, int64, float32, float64, string, []byte:
			propA[k] = v
		case []string:
			switch v.(type) {
			case []string:
				for _, vv := range v.([]string) {
					propA[k] = append(propA[k].([]string), vv)
				}
			default:
				propA[k] = v
			}
		case []int:
			switch v.(type) {
			case []int:
				for _, vv := range v.([]int) {
					propA[k] = append(propA[k].([]int), vv)
				}
			default:
				propA[k] = v
			}
		case map[string]interface{}:
			switch v.(type) {
			case map[string]interface{}:
				obj.CopiedOverMaps(propA[k].(map[string]interface{}), v.(map[string]interface{}))
			default:
				propA[k] = v
			}
		}
	}
}
