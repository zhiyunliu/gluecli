package model

import (
	"strings"

	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
)

type TmplTable struct {
	DbType  string    //数据库类型
	Name    string    //表名
	Desc    string    //表描述
	ExtInfo string    //扩展信息
	Cols    *TmplCols //原始列
	Indexs  *TmplIdxs //索引

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

func (tc *TmplTable) GetPks() *TmplIdx {
	idxs := tc.GetIdxs()
	for _, v := range idxs.Map {
		if v.IdxType == indextype.PK {
			return v
		}
	}
	return nil
}

func (tc *TmplTable) GetIdxs() *TmplIdxs {
	if tc.Indexs != nil {
		return tc.Indexs
	}
	tblIdx := &TmplIdxs{
		Table: tc,
		Map:   map[string]*TmplIdx{},
	}
	tc.Indexs = tblIdx

	for _, col := range tc.Cols.Cols {
		//处理PK
		for k, v := range col.GetPK() {
			idx, ok := tblIdx.Map[k]
			if !ok {
				idx = &TmplIdx{Table: tc, IdxType: indextype.PK, Name: k}
				tblIdx.Map[k] = idx
			}
			idx.Cols = append(idx.Cols, tmplIdxCol{ColName: col.ColName, Sort: v})
		}

		//处理Idx
		for k, v := range col.GetIdxs() {
			idx, ok := tblIdx.Map[k]
			if !ok {
				idx = &TmplIdx{Table: tc, IdxType: indextype.Idx, Name: k}
				tblIdx.Map[k] = idx
			}
			idx.Cols = append(idx.Cols, tmplIdxCol{ColName: col.ColName, Sort: v})
		}

		//处理UNQ
		for k, v := range col.GetUnq() {
			idx, ok := tblIdx.Map[k]
			if !ok {
				idx = &TmplIdx{Table: tc, IdxType: indextype.Unq, Name: k}
				tblIdx.Map[k] = idx
			}
			idx.Cols = append(idx.Cols, tmplIdxCol{ColName: col.ColName, Sort: v})
		}
	}
	tblIdx.Sort()
	return tblIdx
}
