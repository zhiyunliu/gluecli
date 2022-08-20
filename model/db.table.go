package model

type DbTable struct {
	Name    string
	Cols    []DbColInfo
	Idxs    []DbIdx
	Comment string
}
