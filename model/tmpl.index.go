package model

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/zhiyunliu/gluecli/consts/enums/difftype"
	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
)

type TmplIdxs struct {
	Table  *TmplTable
	idxMap map[string]*TmplIdx
}

func NewTmplIdxs(tbl *TmplTable) *TmplIdxs {
	return &TmplIdxs{
		Table:  tbl,
		idxMap: map[string]*TmplIdx{},
	}
}

func (t *TmplIdxs) Append(idx *TmplIdx) error {
	if _, ok := t.idxMap[idx.Name]; ok {
		return fmt.Errorf("[%s]已存在索引[%s]", t.Table.Name, idx.Name)
	}
	idx.Table = t.Table
	t.idxMap[idx.Name] = idx
	return nil
}

func (t *TmplIdxs) GetIdx(name string) (idx *TmplIdx, ok bool) {
	idx, ok = t.idxMap[name]
	return
}

func (t *TmplIdxs) GetIdxTypeList(idxType string) []*TmplIdx {
	result := make(tmplIdxList, 0)
	for _, v := range t.idxMap {
		if strings.EqualFold(v.IdxType, idxType) {
			result = append(result, v)
		}
	}
	sort.Sort(result)
	return result
}

func (t *TmplIdxs) GetIdxList() []*TmplIdx {
	list := make([]*TmplIdx, 0, len(t.idxMap))
	list = append(list, t.GetIdxTypeList(indextype.PK)...)
	list = append(list, t.GetIdxTypeList(indextype.Idx)...)
	list = append(list, t.GetIdxTypeList(indextype.Unq)...)
	return list
}

func (t *TmplIdxs) SortCols() {
	for _, v := range t.idxMap {
		sort.Sort(v.Cols)
	}
}

func (t *TmplIdxs) Diff(target *TmplIdxs) []*TmplIdx {
	tempSource := t.idxMap
	diff := make([]*TmplIdx, 0)

	//减少,索引要先处理删除
	for name, idx := range target.idxMap {
		if _, ok := tempSource[name]; !ok {
			idx.Operation = difftype.Delete
			diff = append(diff, idx)
		}
	}

	//新增
	for name, index := range tempSource {
		if _, ok := target.idxMap[name]; !ok {
			index.Operation = difftype.Insert
			diff = append(diff, index)
			delete(tempSource, name)
		}
	}

	//变动
	for name, sourceIndex := range tempSource {
		if !sourceIndex.Equal(target.idxMap[name]) {
			sourceIndex.Operation = difftype.Modify
			diff = append(diff, sourceIndex)
		}
	}

	return diff
}

type tmplIdxList []*TmplIdx

func (list tmplIdxList) Len() int {
	return len(list)
}

func (list tmplIdxList) Less(i, j int) bool {
	return list[i].Name < list[j].Name
}

// Swap swaps the elements with indexes i and j.
func (list tmplIdxList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

type tmplIdxColList []tmplIdxCol
type TmplIdx struct {
	Table   *TmplTable
	IdxType string
	Name    string
	Cols    tmplIdxColList
	//**************
	Operation difftype.Operation
}

func (s *TmplIdx) Equal(t *TmplIdx) bool {
	if s.Name != t.Name {
		return false
	}
	if !strings.EqualFold(s.IdxType, t.IdxType) {
		return false
	}
	return reflect.DeepEqual(s.Cols, t.Cols)
}

type tmplIdxCol struct {
	ColName string
	Sort    int
}

func (list tmplIdxColList) Len() int {
	return len(list)
}

func (list tmplIdxColList) Less(i, j int) bool {
	return list[i].Sort < list[j].Sort
}

// Swap swaps the elements with indexes i and j.
func (list tmplIdxColList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
