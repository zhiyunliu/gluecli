package model

type DbColInfo struct {
	ColName string
	ColType string
	IsNull  bool
	Default interface{}
	Comment string
}
