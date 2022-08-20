package difftype

// Operation defines the operation of a diff item.
type DBColPart int8

//go:generate stringer -type=Operation -trimprefix=Diff

const (
	//列属性
	ColProperty DBColPart = 1
	//列默认值
	ColDefault DBColPart = 2
	//列描述
	ColComment DBColPart = 3
)
