package funcs

import (
	"reflect"
	"text/template"
)

var BaseFuncs = template.FuncMap{}

func init() {

	BaseFuncs["length"] = func(r interface{}) int {
		v := reflect.ValueOf(r)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Slice || v.Kind() == reflect.Array || v.Kind() == reflect.Map {
			return v.Len()
		}
		return 0
	}

}
