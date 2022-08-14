package model

import "strings"

type TmplTable struct {
	DbType  string    //数据库类型
	Name    string    //表名
	Desc    string    //表描述
	ExtInfo string    //扩展信息
	Cols    *TmplCols //原始行
	Indexs  *TmplIdxs //序列

	//********************
	DropTable bool //生成删除语句
	Exclude   bool //排除生成sql

}

//NewTable 创建表
func NewTmplTable(name, desc, extinfo string) *TmplTable {
	return &TmplTable{
		Name:    strings.TrimLeft(name, "!"),
		Desc:    desc,
		Cols:    &TmplCols{},
		Indexs:  &TmplIdxs{},
		ExtInfo: extinfo,
		Exclude: strings.Contains(name, "!"),
	}
}

func (t *TmplTable) AddCol(col *TmplCol) error {
	t.Cols.Cols = append(t.Cols.Cols, col)
	return nil
}

func (t *TmplTable) GetCol(idx int) *TmplCol {
	return t.Cols.Cols[idx]
}
