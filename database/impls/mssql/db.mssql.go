package mssql

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/zhiyunliu/gluecli/database/define"
	"github.com/zhiyunliu/gluecli/model"
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

func (db *dbMssql) GetDbInfo(args ...interface{}) (dbInfo *model.TmplTableList, err error) {

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
	return strings.ReplaceAll(buff.String(), "\r\n\r\n", "\r\n"), nil
}
func (db *dbMssql) Diff(tbl *model.TmplTable) (content string, err error) {
	var tmpl = template.New("table").Funcs(funcMap)
	np, err := tmpl.Parse(TmplDiffSQLModify)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, tbl); err != nil {
		return "", err
	}
	return strings.ReplaceAll(buff.String(), "\r\n\r\n", "\r\n"), nil
}
