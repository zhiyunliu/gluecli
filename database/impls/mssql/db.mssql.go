package mssql

import (
	"github.com/zhiyunliu/gluecli/database/define"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/gluecli/objecttypes"
)

const (
	DbType = "mssql"
)

func init() {

	define.Registry(&dbMssql{})
}

type dbMssql struct{}

func (db *dbMssql) DbType() string {
	return DbType
}

func (db *dbMssql) GetDbInfo(args ...interface{}) (dbInfo *objecttypes.DbInfo, err error) {
	return
}

func (db *dbMssql) BuildScheme(tbl *model.TmplTable) (content string, err error) {
	return
}
func (db *dbMssql) Diff(tbl *model.TmplTable) (content string, err error) {
	return
}
