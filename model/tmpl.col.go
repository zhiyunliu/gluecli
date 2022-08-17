package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
)

var (
	pkPattern  = regexp.MustCompile(`PK(\(\w+(,\d+)?\))?`)
	idxPattern = regexp.MustCompile(`IDX(\(\w+(,\d+)?\))?`)
	unqPattern = regexp.MustCompile(`UNQ(\(\w+(,\d+)?\))?`)
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
	parties := unqPattern.FindAllString(c.Condition, -1)
	result := map[string]int{}
	for i := range parties {
		name, idx := c.splitIdx(parties[i], indextype.Unq)
		result[name] = idx
	}
	return result
}

/**
//PK
//PK(pk_xxx_df)
//PK(pk_xxx_df,1)
*/
func (c *TmplCol) splitIdx(val, idxType string) (name string, idx int) {
	idxType = strings.ToLower(idxType)
	val = strings.ToLower(val)
	val = strings.TrimSpace(val)

	if strings.EqualFold(val, idxType) {
		return fmt.Sprintf("%s_%s", idxType, c.Table.Name), 1
	}
	val = strings.TrimPrefix(val, idxType)
	val = strings.TrimPrefix(val, "(")
	val = strings.TrimSuffix(val, ")")
	// PK(pk_xxx_df)
	if !strings.Contains(val, ",") {
		return val, 1
	}
	ps := strings.SplitN(val, ",", 2)
	name = ps[0]
	tmpidx, _ := strconv.ParseInt(ps[1], 10, 32)
	idx = int(tmpidx)
	return

}
