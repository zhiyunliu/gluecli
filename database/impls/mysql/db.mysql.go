package mysql

import (
	"github.com/zhiyunliu/gluecli/database/define"
	"github.com/zhiyunliu/gluecli/model"
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

func (db *dbMysql) GetDbInfo(args ...interface{}) (dbInfo *model.DbInfo, err error) {
	return
}
func (db *dbMysql) BuildScheme(tbl *model.TmplTable) (content string, err error) {
	return
}
func (db *dbMysql) Diff(tbl *model.TmplTable) (content string, err error) {
	return
}
