package mssql

import (
	"bytes"
	"strings"
	"text/template"

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
	var tmpl = template.New("table").Funcs(funcMap)
	np, err := tmpl.Parse(TmplCreateTable)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, tbl); err != nil {
		return "", err
	}
	return strings.Replace(strings.Replace(buff.String(), "{###}", "`", -1), "&#39;", "'", -1), nil
}
func (db *dbMssql) Diff(tbl *model.TmplTable) (content string, err error) {
	return
}
