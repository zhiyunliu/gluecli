package indextype

type DbIndexType string

const (
	//主键
	PK DbIndexType = "pk"
	//唯一
	Unq DbIndexType = "unq"
	//普通
	Idx DbIndexType = "idx"
)
