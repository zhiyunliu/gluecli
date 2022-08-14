package diffoeration

// Operation defines the operation of a diff item.
type Operation int8

//go:generate stringer -type=Operation -trimprefix=Diff

const (
	//删除
	Delete Operation = -1
	//新增
	Insert Operation = 1
	//相等
	Equal Operation = 0
	//修改
	Modify Operation = 2
)
