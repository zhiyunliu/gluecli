package objecttypes

type DbInfo struct {
	DbType  DbType
	Tables  []DbTable
	ConnStr string
}

type DbTable struct {
	TableName string
	Cols      []DbColInfo
	Idxs      []DbIdx
	Comment   string
}

type DbColInfo struct {
	ColName string
	ColType string
	IsNull  bool
	Default interface{}
	Comment string
}

type DbIdx struct {
	TableName string
	IdxType   string
	Name      string
	Cols      []string
}
