package model

import (
	"regexp"
	"strings"

	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
)

var (
	pkPattern  = regexp.MustCompile(`PK[\(\w+[,\d+]?\)]?`)
	idxPattern = regexp.MustCompile(`IDX\(\w+[,\d+]?\)`)
	unqPattern = regexp.MustCompile(`UNQ\(\w+[,\d+]?\)`)
)

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
	Default    string
	Comment    string
	Condition  string
}

func (c *TmplCol) HasPk() bool {
	parties := strings.Split(c.Condition, ",")
	for i := range parties {
		if strings.EqualFold(strings.TrimSpace(parties[i]), indextype.PK) {
			return true
		}
	}
	return false
}

func (c *TmplCol) GetPK() map[string]int {
	parties := pkPattern.FindAllString(c.Condition, -1)

	result := map[string]int{}

	for i := range parties {
		name, idx := c.splitIdx(parties[i], indextype.PK)
		result[name] = idx
	}

	return result
}

func (c *TmplCol) GetIdxs() map[string]int {
	parties := idxPattern.FindAllString(c.Condition, -1)

	result := map[string]int{}

	for i := range parties {
		name, idx := c.splitIdx(parties[i], indextype.Idx)
		result[name] = idx
	}

	return result
}

func (c *TmplCol) GetUnq() map[string]int {
	//UNQ(name,1)
	parties := idxPattern.FindAllString(c.Condition, -1)
	result := map[string]int{}
	for i := range parties {
		name, idx := c.splitIdx(parties[i], indextype.Unq)
		result[name] = idx
	}
	return result
}

func (c *TmplCol) splitIdx(val, idxType string) (name string, idx int) {
	//PK(pk_xxx_df,1) | PK
	//IDX(idx_xxx_df,1)
	//UNQ(unq_xxx_df,1)
	//return unq_xxx_df,1
	return
}
