package model

import (
	"strings"

	"github.com/zhiyunliu/gluecli/consts/enums/difftype"
	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
)

type TmplTable struct {
	DbType  string    //数据库类型
	Name    string    //表名
	Desc    string    //表描述
	ExtInfo string    //扩展信息
	Cols    *TmplCols //原始列
	indexs  *TmplIdxs //索引

	//********************
	DropTable bool //生成删除语句
	Exclude   bool //排除生成sql

	//********************
	Operation difftype.Operation
	DiffCols  []*TmplCol
	DiffIdxs  []*TmplIdx
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
	col.Table = t
	t.Cols.Cols = append(t.Cols.Cols, col)
	return nil
}

func (t *TmplTable) GetCol(idx int) *TmplCol {
	return t.Cols.Cols[idx]
}

func (tc *TmplTable) GetPks() *TmplIdx {
	idxs := tc.GetIdxs()
	for _, v := range idxs.GetIdxTypeList(indextype.PK) {
		return v
	}
	return nil
}

func (tc *TmplTable) GetIdxs() *TmplIdxs {
	if tc.indexs != nil {
		return tc.indexs
	}
	tblIdx := NewTmplIdxs(tc)
	tc.indexs = tblIdx

	for _, col := range tc.Cols.Cols {
		//处理PK
		for k, v := range col.GetPK() {
			idx, ok := tblIdx.GetIdx(k)
			if !ok {
				idx = &TmplIdx{IdxType: indextype.PK, Name: k}
				tblIdx.Append(idx)
			}
			idx.Cols = append(idx.Cols, tmplIdxCol{ColName: col.ColName, Sort: v})
		}

		//处理Idx
		for k, v := range col.GetIdxs() {
			idx, ok := tblIdx.GetIdx(k)
			if !ok {
				idx = &TmplIdx{IdxType: indextype.Idx, Name: k}
				tblIdx.Append(idx)
			}
			idx.Cols = append(idx.Cols, tmplIdxCol{ColName: col.ColName, Sort: v})
		}

		//处理UNQ
		for k, v := range col.GetUnq() {
			idx, ok := tblIdx.GetIdx(k)
			if !ok {
				idx = &TmplIdx{IdxType: indextype.Unq, Name: k}
				tblIdx.Append(idx)
			}
			idx.Cols = append(idx.Cols, tmplIdxCol{ColName: col.ColName, Sort: v})
		}
	}
	tblIdx.SortCols()

	return tblIdx
}
