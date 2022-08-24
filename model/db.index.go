package model

type DbIdx struct {
	TableName string
	IdxType   string
	Name      string
	Cols      []string
}
