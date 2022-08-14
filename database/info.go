package database

import (
	"github.com/zhiyunliu/gluecli/database/define"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/gluecli/objecttypes"

	_ "github.com/zhiyunliu/gluecli/database/impls/mssql"
	_ "github.com/zhiyunliu/gluecli/database/impls/mysql"
	_ "github.com/zhiyunliu/gluecli/database/impls/oracle"
)

func GetDbInfo(dbType string, args ...interface{}) (info *objecttypes.DbInfo, err error) {
	dbImpl, err := define.Load(dbType)
	if err != nil {
		return
	}
	info, err = dbImpl.GetDbInfo(args)
	return
}

func BuildScheme(dbType string, tbl *model.TmplTable) (content string, err error) {
	dbImpl, err := define.Load(dbType)
	if err != nil {
		return
	}
	content, err = dbImpl.BuildScheme(tbl)
	return
}

func Diff(dbType string, tbl *model.TmplTable) (content string, err error) {
	dbImpl, err := define.Load(dbType)
	if err != nil {
		return
	}
	content, err = dbImpl.Diff(tbl)
	return
}
