package markdown

import (
	"text/template"

	"github.com/zhiyunliu/gluecli/funcs"
	"github.com/zhiyunliu/gluecli/model"
)

func copyBasicFuncs(funcMap template.FuncMap) {
	for k, v := range funcs.BaseFuncs {
		funcMap[k] = v
	}

}

func getfuncs(dbtype string) map[string]interface{} {
	funcMap := template.FuncMap{}
	copyBasicFuncs(funcMap)
	fillFuncMap(dbtype, funcMap)
	return funcMap
}

func fillFuncMap(dbtype string, funcMap template.FuncMap) {
	funcMap["dbcolType"] = func(col *model.DbColInfo) string {
		return ""
	}
	funcMap["isNull"] = func(col *model.DbColInfo) string {
		return ""
	}
	funcMap["colCondition"] = func(col *model.DbColInfo) string {
		return ""
	}

}
