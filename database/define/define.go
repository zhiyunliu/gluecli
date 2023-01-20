package define

import (
	"fmt"
	"sync"

	contribxdb "github.com/zhiyunliu/glue/contrib/xdb"
	"github.com/zhiyunliu/glue/xdb"
	"github.com/zhiyunliu/gluecli/model"
)

var (
	_dbImpls = sync.Map{}
)

type DbImpl interface {
	DbType() string
	GetDbInfo(args ...interface{}) (*model.TmplTableList, error)
	BuildScheme(tbl *model.TmplTable) (content string, err error)
	Diff(tbl *model.TmplTable) (content string, err error)
}

func Registry(db DbImpl) (err error) {
	_, ok := _dbImpls.Load(db.DbType())
	if ok {
		return fmt.Errorf("存在重复的数据库类型:%s,可以调用Deregistry 进行删除后执行", db.DbType())
	}
	_dbImpls.Store(db.DbType(), db)
	return nil
}

func Deregistry(dbType string) {
	//
	_dbImpls.Delete(dbType)
}

func Load(dbType string) (dbImpl DbImpl, err error) {
	tmpImpl, ok := _dbImpls.Load(dbType)
	if !ok {
		err = fmt.Errorf("不存在[%s]的数据库实现", dbType)
		return
	}
	dbImpl, ok = tmpImpl.(DbImpl)
	if !ok {
		err = fmt.Errorf("错误的类型,未实现接口DbImpl")
		return
	}
	return dbImpl, nil
}

func GetDbInstance(proto string, conn string) (instance xdb.IDB, err error) {
	return contribxdb.NewDB(proto, &contribxdb.Config{
		Proto: proto,
		Conn:  conn,
	})
}
