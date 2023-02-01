package mssql

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/zhiyunliu/glue/contrib/xdb/sqlserver"

	"github.com/zhiyunliu/glue/xdb"
	"github.com/zhiyunliu/gluecli/database/define"
	"github.com/zhiyunliu/gluecli/database/impls/mssql/sqls"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/golibs/xtypes"
)

const (
	_Proto = "sqlserver"
)

type CaclLenCallback func(col xtypes.XMap) (int, int)
type dbIdx struct {
	idxtype string
	colsort int
}

func doubleCol(col xtypes.XMap) (int, int) {
	v1, _ := col.GetInt("length")
	v2, _ := col.GetInt("decimal")
	return v1, v2
}

var (
	CalcColType = map[string]CaclLenCallback{
		"date":          func(col xtypes.XMap) (int, int) { return 0, 0 },
		"datetime":      func(col xtypes.XMap) (int, int) { return 0, 0 },
		"smalldatetime": func(col xtypes.XMap) (int, int) { return 0, 0 },
		"timestamp":     func(col xtypes.XMap) (int, int) { return 0, 0 },
		"bigint":        func(col xtypes.XMap) (int, int) { return 0, 0 },
		"int":           func(col xtypes.XMap) (int, int) { return 0, 0 },
		"smallint":      func(col xtypes.XMap) (int, int) { return 0, 0 },
		"tinyint":       func(col xtypes.XMap) (int, int) { return 0, 0 },
		"decimal":       func(col xtypes.XMap) (int, int) { return doubleCol(col) },
		"numeric":       func(col xtypes.XMap) (int, int) { return doubleCol(col) },
		"smallmoney":    func(col xtypes.XMap) (int, int) { return 0, 0 },
		"money":         func(col xtypes.XMap) (int, int) { return 0, 0 },
	}
)

func (db *dbMssql) GetDbInfo(args ...interface{}) (dbInfo *model.TmplTableList, err error) {
	dbObj, err := define.GetDbInstance(_Proto, fmt.Sprint(args[0]))
	if err != nil {
		return
	}

	dbInfo = model.NewTableList()
	ctx := context.Background()

	tableList, err := db.queryTableList(ctx, dbObj, fmt.Sprint(args[1]))
	if err != nil {
		return
	}

	var cols xtypes.XMaps
	var idxs xtypes.XMaps

	for _, tbl := range tableList {
		tmpTbl := model.NewTmplTable(tbl.GetString("name"), tbl.GetString("value"), "")
		tmpTbl.DbType = db.DbType()
		if err = dbInfo.Append(tmpTbl); err != nil {
			return
		}
		cols, err = db.queryTableCols(ctx, dbObj, tbl.GetString("object_id"))
		if err != nil {
			return
		}
		idxs, err = db.queryTableIdxs(ctx, dbObj, tbl.GetString("object_id"))
		if err != nil {
			return
		}
		newIdxs := db.rebuildIdxs(idxs)
		for i, col := range cols {

			colLen, decialLen := db.calcLen(col)

			tplCol := &model.TmplCol{
				LineNum:    i,
				Table:      tmpTbl,
				ColName:    col.GetString("col_name"),
				ColType:    db.buildColType(col),
				ColLen:     colLen,
				DecimalLen: decialLen,
				IsNull:     col.GetString("isnullable"),
				Default:    db.calcDefaultVal(col),
				Comment:    col.GetString("comments"),
				Condition:  db.calcColCondition(col, newIdxs),
			}
			tmpTbl.AddCol(tplCol)
		}
	}

	return
}

func (db *dbMssql) queryTableList(ctx context.Context, dbObj xdb.IDB, tableName string) (tableList xtypes.XMaps, err error) {
	tableList, err = dbObj.Query(ctx, sqls.QueryTables, map[string]interface{}{
		"name": tableName,
	})
	if err != nil {
		return
	}
	return
}

func (db *dbMssql) queryTableCols(ctx context.Context, dbObj xdb.IDB, objId string) (colList xtypes.XMaps, err error) {
	colList, err = dbObj.Query(ctx, sqls.QueryTableCols, map[string]interface{}{
		"obj_id": objId,
	})
	if err != nil {
		return
	}
	return
}
func (db *dbMssql) queryTableIdxs(ctx context.Context, dbObj xdb.IDB, objId string) (idxList xtypes.XMaps, err error) {
	idxList, err = dbObj.Query(ctx, sqls.QueryTableIdxs, map[string]interface{}{
		"obj_id": objId,
	})
	if err != nil {
		return
	}
	return
}

func (db *dbMssql) calcLen(col xtypes.XMap) (collen, decimal int) {
	colType := col.GetString("col_type")
	callback, ok := CalcColType[colType]
	if ok {
		return callback(col)
	}
	collen, _ = col.GetInt("length")
	return
}

func (db *dbMssql) calcDefaultVal(col xtypes.XMap) (dftVal string) {
	dftVal = col.GetString("default_val")
	if dftVal == "" {
		return
	}
	spc := '('
	cnt := 0
	for _, c := range dftVal {
		if c == spc {
			cnt++
			continue
		}
		break
	}
	dftLen := len(dftVal)
	dftVal = dftVal[cnt : dftLen-cnt]
	return
}

func (db *dbMssql) buildColType(col xtypes.XMap) (colType string) {
	colType = col.GetString("col_type")
	colLen, decimalLen := db.calcLen(col)
	if colLen <= 0 {
		return colType
	}
	if decimalLen <= 0 {
		return fmt.Sprintf("%s(%d)", colType, colLen)
	}
	return fmt.Sprintf("%s(%d,%d)", colType, colLen, decimalLen)
}

func (db *dbMssql) calcColCondition(col xtypes.XMap, idxs map[string]map[string]dbIdx) (condition string) {
	vals := []string{}
	if col.GetBool("isidentity") {
		vals = append(vals, "SEQ")
	}
	colName := col.GetString("col_name")

	idxVal, ok := idxs[colName]
	if ok {
		for k, v := range idxVal {
			vals = append(vals, fmt.Sprintf("%s(%s,%d)", v.idxtype, k, v.colsort))
		}
	}
	condition = strings.Join(vals, ",")
	return
}

func (db *dbMssql) rebuildIdxs(idxs xtypes.XMaps) (newidxs map[string]map[string]dbIdx) {
	//col_name[idx_name]property
	newidxs = map[string]map[string]dbIdx{}
	if len(idxs) == 0 {
		return
	}
	for _, idx := range idxs {
		colName := idx.GetString("col_name")
		idxName := idx.GetString("idx_name")
		idxMap, ok := newidxs[colName]
		if !ok {
			newidxs[colName] = map[string]dbIdx{}
			idxMap = newidxs[colName]
		}
		sortVal, _ := idx.GetInt("sort_val")
		idxMap[idxName] = dbIdx{
			idxtype: idx.GetString("idx_type"),
			colsort: sortVal,
		}
	}
	return
}
