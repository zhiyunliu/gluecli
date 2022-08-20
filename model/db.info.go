package model

type DbInfo struct {
	DbType  string
	Tables  []DbTable
	ConnStr string
}
