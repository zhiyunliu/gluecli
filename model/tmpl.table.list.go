package model

import "strings"

type TmplTableList struct {
	Tables []*TmplTable
	Map    map[string]bool
}

//FilterTable 过滤行信息(逗号分隔)
func (t *TmplTableList) FilterTable(tableNameList string) {
	if tableNameList == "" {
		return
	}
	kws := strings.Split(tableNameList, ",")
	tbs := make([]*TmplTable, 0, len(kws))
	for _, v := range kws {
		for _, tb := range t.Tables {
			if strings.Contains(tb.Name, v) {
				tbs = append(tbs, tb)
			}
		}
	}
	t.Tables = tbs
}
func (t *TmplTableList) DropTable(dropTable bool) {
	for _, tb := range t.Tables {
		tb.DropTable = dropTable
	}
}

func (t *TmplTableList) Exclude() {
	tbs := make([]*TmplTable, 0, 1)
	for _, tb := range t.Tables {
		if !tb.Exclude {
			tbs = append(tbs, tb)
		}
	}
	t.Tables = tbs
}
func (t *TmplTableList) Diff(n *TmplTableList) *TmplTableList {
	return t
}
