package mssql

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
	"github.com/zhiyunliu/gluecli/funcs"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/golibs/xtransform"
	"github.com/zhiyunliu/golibs/xtypes"
)

var (
	funcMap = template.FuncMap{}

	nullMap = xtypes.XMap{
		"":  "NOT NULL ",
		"N": "NOT NULL ",
		"n": "NOT NULL ",
		"否": "NOT NULL ",
		"Y": "NULL",
		"y": "NULL",
		"是": "NULL",
	}

	defaultMap = xtypes.XMap{
		"sysdate":           "getdate()",
		"current_timestamp": "getdate()",
		"now()":             "getdate()",
	}

	colTypeMap = xtypes.XMap{
		"varchar2": "varchar",
	}
)

func init() {
	for k, v := range funcs.BaseFuncs {
		funcMap[k] = v
	}

	funcMap["dbcolType"] = func(col *model.TmplCol) string {
		colType := strings.ToLower(col.ColType)
		v := colTypeMap.GetString(colType)
		if v == "" {
			v = colType
		}
		partList := []string{}
		if col.ColLen != 0 {
			partList = append(partList, strconv.Itoa(col.ColLen))
		}
		if col.DecimalLen != 0 {
			partList = append(partList, strconv.Itoa(col.DecimalLen))
		}
		partVal := ""
		if len(partList) > 0 {
			partVal = strings.Join(partList, ",")
			partVal = "(" + partVal + ")"
		}
		return colType + partVal
	}

	funcMap["defaultValue"] = func(col *model.TmplCol) string {
		col.Default = strings.TrimSpace(col.Default)

		if strings.EqualFold(col.Default, "") {
			return ""
		}
		col.Default = strings.ToLower(col.Default)
		return defaultMap.GetString(col.Default)
	}

	funcMap["isNull"] = func(col *model.TmplCol) string {
		return nullMap.GetString(strings.TrimSpace(col.IsNull))
	}

	funcMap["generatePK"] = func(tbl *model.TmplTable) string {
		tmpl := `
		CONSTRAINT @{Name} PRIMARY KEY CLUSTERED 
		(
			@{PkCols}
		) `

		pks := tbl.GetPks()

		pkslist := []string{}

		for i := range pks.Cols {
			pkslist = append(pkslist, pks.Cols[i].ColName+" ASC ")
		}

		result := xtransform.Translate(tmpl, map[string]interface{}{
			"Name":   tbl.Name,
			"PkCols": strings.Join(pkslist, ","),
		})

		return result
	}

	funcMap["generateIdx"] = func(tbl *model.TmplTable) string {
		idxtmpl := `
		CREATE @{UNIQUE} NONCLUSTERED INDEX @{idx_name} ON @{table_name}( @{idx_cols} )
		`

		idxs := tbl.GetIdxs()

		list := []string{}

		for k, v := range idxs.Map {
			param := map[string]interface{}{
				"idx_name":   k,
				"table_name": tbl.Name,
			}
			if strings.EqualFold(v.IdxType, indextype.Unq) {
				param["UNIQUE"] = "UNIQUE"
			}
			cols := make([]string, len(v.Cols))
			for i, c := range v.Cols {
				cols[i] = fmt.Sprintf("%s ASC ", c.ColName)
			}
			param["idx_cols"] = strings.Join(cols, ",")

			list = append(list, xtransform.Translate(idxtmpl, param))
		}

		return strings.Join(list, "\r\n")
	}

	funcMap["generateComment"] = func(tbl *model.TmplTable) string {
		tmpl := `EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'@{Comment}' , @level0type=N'SCHEMA',@level0name=N'dbo', @level1type=N'TABLE',@level1name=N'@{TableName}', @level2type=N'COLUMN',@level2name=N'@{ColName}'`
		list := []string{}
		for _, col := range tbl.Cols.Cols {
			param := map[string]interface{}{
				"Comment":   col.Comment,
				"TableName": tbl.Name,
				"ColName":   col.ColName,
			}
			list = append(list, xtransform.Translate(tmpl, param))
		}
		return strings.Join(list, "\r\n")
	}
}
