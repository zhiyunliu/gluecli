package model

import (
	"reflect"
	"sort"
	"strings"

	"github.com/zhiyunliu/gluecli/consts/enums/diffoeration"
)

type TmplIdxs struct {
	Table *TmplTable
	Map   map[string]*TmplIdx
}

func (t *TmplIdxs) Sort() {
	for _, v := range t.Map {
		sort.Sort(v.Cols)
	}
}

func (t *TmplIdxs) Diff(target *TmplIdxs) []*TmplIdx {
	tempSource := t.Map
	diff := make([]*TmplIdx, 0)

	//减少,索引要先处理删除
	for name, idx := range target.Map {
		if _, ok := tempSource[name]; !ok {
			idx.Operation = diffoeration.Delete
			diff = append(diff, idx)
		}
	}

	//新增
	for name, index := range tempSource {
		if _, ok := target.Map[name]; !ok {
			index.Operation = diffoeration.Insert
			diff = append(diff, index)
			delete(tempSource, name)
		}
	}

	//变动
	for name, sourceIndex := range tempSource {
		if !sourceIndex.Equal(target.Map[name]) {
			sourceIndex.Operation = diffoeration.Modify
			diff = append(diff, sourceIndex)
		}
	}

	return diff
}

type tmplIdxColList []tmplIdxCol
type TmplIdx struct {
	Table   *TmplTable
	IdxType string
	Name    string
	Cols    tmplIdxColList
	//**************
	Operation diffoeration.Operation
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
