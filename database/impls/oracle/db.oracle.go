package oracle

import (
	"github.com/zhiyunliu/gluecli/database/define"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/gluecli/objecttypes"
)

const (
	DbType = "oracle"
)

func init() {

	define.Registry(&dbOracle{})
}

type dbOracle struct{}

func (db *dbOracle) DbType() string {
	return DbType
}

func (db *dbOracle) GetDbInfo(args ...interface{}) (dbInfo *objecttypes.DbInfo, err error) {
	return
}
func (db *dbOracle) BuildScheme(tbl *model.TmplTable) (content string, err error) {
	return
}
func (db *dbOracle) Diff(tbl *model.TmplTable) (content string, err error) {
	return
}
