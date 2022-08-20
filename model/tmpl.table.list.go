package model

import (
	"fmt"
	"strings"

	"github.com/zhiyunliu/gluecli/consts/enums/diffoeration"
)

type TmplTableList struct {
	Tables []*TmplTable
	tblMap map[string]bool
}

func NewTableList() *TmplTableList {
	return &TmplTableList{
		Tables: make([]*TmplTable, 0),
		tblMap: make(map[string]bool),
	}
}

func (t *TmplTableList) Append(tbl *TmplTable) (err error) {
	if _, ok := t.tblMap[tbl.Name]; ok {
		return fmt.Errorf("存在相同的表名：%s", tbl.Name)
	}
	t.Tables = append(t.Tables, tbl)
	return nil
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

func (source *TmplTableList) getTableMap() map[string]*TmplTable {
	m := make(map[string]*TmplTable, len(source.Tables))
	for _, t := range source.Tables {
		m[t.Name] = t
	}
	return m
}

func (t *TmplTableList) Diff(dest *TmplTableList) *TmplTableList {

	sourceM := t.getTableMap()
	targetM := dest.getTableMap()

	diff := &TmplTableList{
		Tables: make([]*TmplTable, 0),
	}

	//新增
	for tname, table := range sourceM {
		if _, ok := targetM[tname]; !ok {
			table.Operation = diffoeration.Insert
			diff.Tables = append(diff.Tables, table)
			delete(sourceM, tname)
		}
	}

	//减少
	for tname, table := range targetM {
		if _, ok := sourceM[tname]; !ok {
			table.Operation = diffoeration.Delete
			diff.Tables = append(diff.Tables, table)
			delete(targetM, tname)
		}
	}

	//变动
	for name, sourceTable := range sourceM {
		targetTable := targetM[name]
		//字段变化
		diffCols := sourceTable.Cols.Diff(targetTable.Cols)
		if len(diffCols) > 0 {
			sourceTable.DiffCols = diffCols
		}
		//索引变化
		sidx := sourceTable.GetIdxs()
		didx := targetTable.GetIdxs()
		diffIdxs := sidx.Diff(didx)
		if len(diffIdxs) > 0 {
			sourceTable.DiffIdxs = diffIdxs
		}
		if len(diffCols) > 0 || len(diffIdxs) > 0 {
			sourceTable.Operation = diffoeration.Modify
			diff.Tables = append(diff.Tables, sourceTable)
		}
	}
	return diff
}
