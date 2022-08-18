package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/zhiyunliu/gluecli/consts/enums/diffoeration"
	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
)

var (
	seqPattern = regexp.MustCompile(`SEQ(\((\d+(,\d+)?)\))?`)
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

func (tc *TmplCols) getColMap() map[string]*TmplCol {
	m := make(map[string]*TmplCol, len(tc.Cols))
	for _, c := range tc.Cols {
		m[c.ColName] = c
	}
	return m
}

func (tc *TmplCols) Diff(dest *TmplCols) []*TmplCol {

	sourceM := tc.getColMap()
	targetM := dest.getColMap()

	diff := make([]*TmplCol, 0)

	//新增
	for name, col := range sourceM {
		if _, ok := targetM[name]; !ok {
			col.Operation = diffoeration.Insert
			diff = append(diff, col)
			delete(sourceM, name)
		}
	}

	//减少
	for name, col := range targetM {
		if _, ok := sourceM[name]; !ok {
			col.Operation = diffoeration.Delete
			diff = append(diff, col)
			delete(targetM, name)
		}
	}

	//变动
	for name, scol := range sourceM {
		if !scol.Equal(targetM[name]) {
			scol.Operation = diffoeration.Modify
			diff = append(diff, scol)
		}
	}

	return diff

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

	//*****************
	Operation diffoeration.Operation
}

func (c *TmplCol) Equal(t *TmplCol) bool {
	return strings.EqualFold(c.ColName, t.ColName) &&
		strings.EqualFold(c.ColType, t.ColType) &&
		strings.EqualFold(c.Default, t.Default) &&
		strings.EqualFold(c.IsNull, t.IsNull) &&
		strings.EqualFold(c.Comment, t.Comment) &&
		c.ColLen == t.ColLen &&
		c.DecimalLen == t.DecimalLen
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

func (c *TmplCol) GetSeq() *KV {
	parties := seqPattern.FindAllString(c.Condition, -1)
	for i := range parties {
		vals := seqPattern.FindStringSubmatch(parties[i])
		tmpV := vals[2]
		if tmpV == "" {
			return &KV{K: "1", V: "1"}
		}
		if !strings.Contains(tmpV, ",") {
			return &KV{K: tmpV, V: "1"}
		}
		ps := strings.SplitN(tmpV, ",", 2)
		return &KV{K: ps[0], V: ps[1]}
	}

	return nil
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
