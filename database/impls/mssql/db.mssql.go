package mssql

import (
	"bytes"
	"strings"
	"text/template"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/zhiyunliu/gluecli/consts/enums/difftype"
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
	txtTmpl := TmplDiffSQLModify
	if tbl.Operation == difftype.Delete {
		txtTmpl = TmplDropTable
	}
	if tbl.Operation == difftype.Insert {
		txtTmpl = TmplCreateTable
	}

	np, err := tmpl.Parse(txtTmpl)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, tbl); err != nil {
		return "", err
	}
	return strings.ReplaceAll(buff.String(), "\r\n\r\n", "\r\n"), nil
}
