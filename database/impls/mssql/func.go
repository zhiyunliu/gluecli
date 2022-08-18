package mssql

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/zhiyunliu/gluecli/consts/enums/indextype"
	"github.com/zhiyunliu/gluecli/funcs"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/golibs/xtransform"
	"github.com/zhiyunliu/golibs/xtypes"
)

type ColTypeCallback func(*model.TmplCol) string

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

	colTypeMap = map[string]string{
		`^varchar\((\d+)\)$`:      "varchar(*)",
		`^varchar2\((\d+)\)$`:     "varchar(*)",
		`^nvarchar2\((\d+)\)$`:    "nvarchar(*)",
		`^number\((\d+),(\d+)\)$`: "decimal(*)",
		`^date$`:                  "datetime",
		`^datetime$`:              "datetime",
		`^timestamp$`:             "datetime",
		`^decimal$`:               "decimal",
		`^float$`:                 "float",
		`^int$`:                   "int",
		`^number\([1-2]{1}\)$`:    "tinyint",
		`^number\([3-9]{1}\)$`:    "int",
		`^number\(10\)$`:          "int",
		`^number\(1[1-9]{1}\)$`:   "bigint",
		`^number\(2[0-9]{1}\)$`:   "bigint",
		`^string$`:                "tinytext",
		`^text$`:                  "text",
		`^longtext$`:              "text",
		`^clob$`:                  "text",
	}
)

func init() {
	for k, v := range funcs.BaseFuncs {
		funcMap[k] = v
	}

	funcMap["dbcolType"] = func(col *model.TmplCol) string {
		colType := strings.ToLower(col.ColType)
		colType = strings.TrimSpace(colType)
		for regx, v := range colTypeMap {
			reg := regexp.MustCompile(regx)
			if reg.MatchString(colType) {
				if !strings.Contains(v, "*") {
					return v
				}
				value := reg.FindStringSubmatch(colType)
				if len(value) > 1 {
					return strings.Replace(v, "*", strings.Join(value[1:], ","), -1)
				}
				return v
			}
		}
		return colType
	}

	funcMap["seq"] = func(col *model.TmplCol) string {
		seqV := col.GetSeq()
		if seqV == nil {
			return ""
		}
		return fmt.Sprintf(" IDENTITY(%s,%s)", seqV.K, seqV.V)
	}

	funcMap["defaultValue"] = func(col *model.TmplCol) string {
		col.Default = strings.TrimSpace(col.Default)

		if strings.EqualFold(col.Default, "") {
			return ""
		}
		newVal := defaultMap.GetString(col.Default)
		if newVal == "" {
			newVal = col.Default
		}

		switch col.ColType {
		case "varchar":
			fallthrough
		case "nvarchar":
			newVal = fmt.Sprintf("'%s'", newVal)
		default:
		}

		return fmt.Sprintf(" default %s", newVal)
	}

	funcMap["isNull"] = func(col *model.TmplCol) string {
		return nullMap.GetString(strings.TrimSpace(col.IsNull))
	}

	funcMap["generatePK"] = func(tbl *model.TmplTable) string {
		tmpl := `
CONSTRAINT @{Name} PRIMARY KEY CLUSTERED 
(
	@{PkCols}
) 
`

		pks := tbl.GetPks()
		if pks == nil {
			panic(fmt.Errorf("【%s】未设置PK", tbl.Name))
		}

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
		idxtmpl := `CREATE @{UNIQUE} NONCLUSTERED INDEX @{idx_name} ON @{table_name}( @{idx_cols} )	`

		idxs := tbl.GetIdxs()

		list := []string{}

		for k, v := range idxs.Map {
			if strings.EqualFold(v.IdxType, indextype.PK) {
				continue
			}
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
		list := []string{}
		for _, col := range tbl.Cols.Cols {
			r := colComment(col)
			list = append(list, r)
		}
		return strings.Join(list, "\r\n")
	}

	funcMap["colComment"] = colComment

	//idx
	funcMap["isPk"] = func(idx *model.TmplIdx) bool {
		return strings.EqualFold(idx.IdxType, indextype.PK)
	}
	funcMap["isIndex"] = func(idx *model.TmplIdx) bool {
		return strings.EqualFold(idx.IdxType, indextype.Idx)
	}
	funcMap["isUNQ"] = func(idx *model.TmplIdx) bool {
		return strings.EqualFold(idx.IdxType, indextype.Unq)
	}

	funcMap["indexStr"] = func(idx *model.TmplIdx) string {
		return idx.Name
	}

}

func colComment(col *model.TmplCol) string {
	tmpl := `EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'%s' , @level0type=N'SCHEMA',@level0name=N'dbo', @level1type=N'TABLE',@level1name=N'%s', @level2type=N'COLUMN',@level2name=N'%s'`
	return fmt.Sprintf(tmpl, col.Comment, col.Table.Name, col.ColName)
}
