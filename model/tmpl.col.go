package model

type TmplCols struct {
	Cols []*TmplCol
}

func (tc *TmplCols) Count() int {
	return len(tc.Cols)
}

type TmplCol struct {
	LineNum    int
	Table      *TmplTable
	ColName    string
	ColType    string
	ColLen     int
	DecimalLen int
	IsNull     string
	Default    interface{}
	Comment    string
	Condition  string
}
