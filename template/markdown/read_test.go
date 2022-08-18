package markdown

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/bmizerany/assert"
)

func TestRead(t *testing.T) {
	text := `## 一、商户信息

	###  1. 商户信息[ots_merchant_info]
	
	| 字段名      | 类型         | 默认值  | 为空  |     约束     | 描述     |
	| ----------- | ------------ | :-----: | :---: | :----------: | :------- |
	| mer_no      | varchar2(32) |         |  否   | PK,SEQ,RL,DI | 编号     |
	| mer_name    | varchar2(64) |         |  否   |   CRUQL,DN   | 名称     |
	| mer_crop    | varchar2(64) |         |  否   |   CRUQL,DN   | 公司     |
	| mer_type    | number(1)    |         |  否   |   CRUQL,DN   | 类型     |
	| bd_uid      | number(20)   |         |  否   |   CRUQL,DN   | 商务人员 |
	| status      | number(1)    |    0    |  否   |   RUQL,SL    | 状态     |
	| create_time | date         | sysdate |  否   |    RL,DTIME     | 创建时间 |`

	b := bytes.NewBuffer([]byte(text))
	lines := loadLinesFromReader(bufio.NewReader(b))
	assert.Equal(t, 11, len(lines))

	tl := line2TableLine(lines)
	assert.Equal(t, 1, len(tl.Lines))

	tb, err := tableLine2Table(tl)
	assert.Equal(t, nil, err)

	assert.Equal(t, 1, len(tb.Tables))
	assert.Equal(t, 7, tb.Tables[0].Cols.Count())
	assert.Equal(t, "ots_merchant_info", tb.Tables[0].Name)
	assert.Equal(t, "商户信息", tb.Tables[0].Desc)
	assert.Equal(t, "create_time", tb.Tables[0].GetCol(6).ColName)
	assert.Equal(t, "date", tb.Tables[0].GetCol(6).ColType)
	assert.Equal(t, "sysdate", tb.Tables[0].GetCol(6).Default)
	assert.Equal(t, "否", tb.Tables[0].GetCol(6).IsNull)
	assert.Equal(t, "RL,DTIME", tb.Tables[0].GetCol(6).Condition)
	assert.Equal(t, "创建时间", tb.Tables[0].GetCol(6).Comment)
}
