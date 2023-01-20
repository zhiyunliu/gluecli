package mssql

import (
	"context"
	"fmt"

	"github.com/zhiyunliu/glue/xdb"
	"github.com/zhiyunliu/gluecli/database/define"
	"github.com/zhiyunliu/gluecli/database/impls/mssql/sqls"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/golibs/xtypes"
)

const (
	_Proto = "sqlserver"
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
		cols.Append(nil)
		idxs.Append(nil)
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
