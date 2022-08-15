package model

import "sort"

type TmplIdxs struct {
	Table *TmplTable
	Map   map[string]*TmplIdx
}

func (t *TmplIdxs) Sort() {
	for _, v := range t.Map {
		sort.Sort(v.Cols)
	}
}

type tmplIdxColList []tmplIdxCol
type TmplIdx struct {
	Table   *TmplTable
	IdxType string
	Name    string
	Cols    tmplIdxColList
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
