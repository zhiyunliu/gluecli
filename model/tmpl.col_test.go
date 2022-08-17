package model

import (
	"reflect"
	"testing"

	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
)

func TestTmplCol_splitIdx(t *testing.T) {

	tests := []struct {
		name     string
		c        *TmplCol
		val      string
		idxType  string
		wantName string
		wantIdx  int
	}{
		{name: "1.1", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "PK", idxType: indextype.PK, wantName: "pk_test", wantIdx: 1},
		{name: "1.2", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "PK(pk_test1)", idxType: indextype.PK, wantName: "pk_test1", wantIdx: 1},
		{name: "1.3", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "PK(pk_test2,2)", idxType: indextype.PK, wantName: "pk_test2", wantIdx: 2},

		{name: "2.1", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "IDX", idxType: indextype.Idx, wantName: "idx_test", wantIdx: 1},
		{name: "2.2", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "IDX(idx_test1)", idxType: indextype.Idx, wantName: "idx_test1", wantIdx: 1},
		{name: "2.3", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "IDX(idx_test2,2)", idxType: indextype.Idx, wantName: "idx_test2", wantIdx: 2},

		{name: "3.1", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "UNQ", idxType: indextype.Unq, wantName: "unq_test", wantIdx: 1},
		{name: "3.2", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "UNQ(unq_test1)", idxType: indextype.Unq, wantName: "unq_test1", wantIdx: 1},
		{name: "3.3", c: &TmplCol{Table: &TmplTable{Name: "test"}}, val: "UNQ(unq_test2,2)", idxType: indextype.Unq, wantName: "unq_test2", wantIdx: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotIdx := tt.c.splitIdx(tt.val, tt.idxType)
			if gotName != tt.wantName {
				t.Errorf("TmplCol.splitIdx() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotIdx != tt.wantIdx {
				t.Errorf("TmplCol.splitIdx() gotIdx = %v, want %v", gotIdx, tt.wantIdx)
			}
		})
	}
}
func TestTmplCol_GetPK(t *testing.T) {
	tests := []struct {
		name string
		c    *TmplCol
		want map[string]int
	}{
		{name: "1.1.", c: &TmplCol{Condition: "PK", Table: &TmplTable{Name: "test"}}, want: map[string]int{"pk_test": 1}},
		{name: "1.2.", c: &TmplCol{Condition: "PK(pk_aaa_1)"}, want: map[string]int{"pk_aaa_1": 1}},
		{name: "1.3.", c: &TmplCol{Condition: "PK(pk_aaa_2,2)"}, want: map[string]int{"pk_aaa_2": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetPK(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TmplCol.GetPK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmplCol_GetIdxs(t *testing.T) {
	tests := []struct {
		name string
		c    *TmplCol
		want map[string]int
	}{
		{name: "2.1.", c: &TmplCol{Condition: "IDX", Table: &TmplTable{Name: "test"}}, want: map[string]int{"idx_test": 1}},
		{name: "2.2.", c: &TmplCol{Condition: "IDX(idx_aaa_1)"}, want: map[string]int{"idx_aaa_1": 1}},
		{name: "2.3.", c: &TmplCol{Condition: "IDX(idx_aaa_2,2)"}, want: map[string]int{"idx_aaa_2": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetIdxs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TmplCol.GetIdxs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTmplCol_GetUnq(t *testing.T) {
	tests := []struct {
		name string
		c    *TmplCol
		want map[string]int
	}{
		{name: "3.1.", c: &TmplCol{Condition: "UNQ", Table: &TmplTable{Name: "test"}}, want: map[string]int{"unq_test": 1}},
		{name: "3.2.", c: &TmplCol{Condition: "UNQ(unq_aaa_1)"}, want: map[string]int{"unq_aaa_1": 1}},
		{name: "3.3.", c: &TmplCol{Condition: "UNQ(unq_aaa_2,2)"}, want: map[string]int{"unq_aaa_2": 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetUnq(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TmplCol.GetUnq() = %v, want %v", got, tt.want)
			}
		})
	}
}
