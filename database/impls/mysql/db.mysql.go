package mysql

import (
	"github.com/zhiyunliu/gluecli/database/define"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/gluecli/objecttypes"
)

const (
	DbType = "mysql"
)

func init() {
	define.Registry(&dbMysql{})
}

type dbMysql struct{}

func (db *dbMysql) DbType() string {
	return DbType
}

func (db *dbMysql) GetDbInfo(args ...interface{}) (dbInfo *objecttypes.DbInfo, err error) {
	return
}
func (db *dbMysql) BuildScheme(tbl *model.TmplTable) (content string, err error) {
	return
}
func (db *dbMysql) Diff(tbl *model.TmplTable) (content string, err error) {
	return
}
