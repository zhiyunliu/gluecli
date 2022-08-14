package model

type TmplIdxs struct {
	Table *TmplTable
	Map   map[string]*TmplIdx
}
type TmplIdx struct {
	Table   *TmplTable
	IdxType string
	Name    string
	Cols    []string
}
