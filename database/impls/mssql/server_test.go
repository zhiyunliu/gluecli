package mssql

import (
	"reflect"
	"testing"

	"github.com/zhiyunliu/golibs/xtypes"
)

func Test_dbMssql_calcDefaultVal(t *testing.T) {
	db := dbMssql{}

	tests := []struct {
		name       string
		col        xtypes.XMap
		wantDftVal string
	}{
		{name: "0.", col: xtypes.XMap{"default_val": ""}, wantDftVal: ""},
		{name: "1.", col: xtypes.XMap{"default_val": "(getdate())"}, wantDftVal: "getdate()"},
		{name: "2.", col: xtypes.XMap{"default_val": "((0))"}, wantDftVal: "0"},
		{name: "3.", col: xtypes.XMap{"default_val": "('aabbcc')"}, wantDftVal: "'aabbcc'"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDftVal := db.calcDefaultVal(tt.col); gotDftVal != tt.wantDftVal {
				t.Errorf("dbMssql.calcDefaultVal() = %v, want %v", gotDftVal, tt.wantDftVal)
			}
		})
	}
}

func Test_dbMssql_rebuildIdxs(t *testing.T) {
	db := dbMssql{}

	tests := []struct {
		name        string
		idxs        xtypes.XMaps
		wantNewidxs map[string]map[string]dbIdx
	}{
		{name: "0.", idxs: xtypes.XMaps{}, wantNewidxs: map[string]map[string]dbIdx{}},
		{name: "1.", idxs: xtypes.XMaps{{"sort_val": 1, "idx_type": "PK", "idx_name": "pk_a", "col_name": "a"}}, wantNewidxs: map[string]map[string]dbIdx{"a": {"pk_a": {idxtype: "PK", colsort: 1}}}},
		{name: "2.", idxs: xtypes.XMaps{{"sort_val": 1, "idx_type": "IDX", "idx_name": "idx_ab", "col_name": "a"}, {"sort_val": 2, "idx_type": "IDX", "idx_name": "idx_ab", "col_name": "b"}}, wantNewidxs: map[string]map[string]dbIdx{"a": {"idx_ab": {idxtype: "IDX", colsort: 1}}, "b": {"idx_ab": {idxtype: "IDX", colsort: 2}}}},
		{name: "3.", idxs: xtypes.XMaps{{"sort_val": 1, "idx_type": "IDX", "idx_name": "idx_ab", "col_name": "a"}, {"sort_val": 1, "idx_type": "UNQ", "idx_name": "unq_ab", "col_name": "b"}}, wantNewidxs: map[string]map[string]dbIdx{"a": {"idx_ab": {idxtype: "IDX", colsort: 1}}, "b": {"unq_ab": {idxtype: "UNQ", colsort: 1}}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNewidxs := db.rebuildIdxs(tt.idxs); !reflect.DeepEqual(gotNewidxs, tt.wantNewidxs) {
				t.Errorf("dbMssql.rebuildIdxs() = %v, want %v", gotNewidxs, tt.wantNewidxs)
			}
		})
	}
}
