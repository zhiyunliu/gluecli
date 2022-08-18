package mssql

const TmplCreateTable = `
{{$count:=.Cols.Cols|length -}}

{{- if .DropTable}}
	
if exists (select * from sysobjects where id = object_id('{{.Name}}') and OBJECTPROPERTY(id, 'IsUserTable') = 1)
	drop table {{.Name}}
go 

{{end -}}
	
CREATE TABLE {{.Name}} (
		{{range $i,$c:=.Cols.Cols -}}
		{{$c.ColName}} {{$c|dbcolType}} {{$c|seq}} {{$c|isNull}}  {{$c|defaultValue}},
		{{end -}}
	{{.|generatePK}}
		
)

{{.|generateIdx}}


go 
{{.|generateComment}}
  
go 
  `

const TmplDiffSQLModify = `

{{- range $i,$c:=.DiffCols}}
{{- if (eq $c.Operation -1)}}
-- 删除字段 {{$c.ColName}} 
ALTER TABLE {{$.Name}} drop COLUMN {{$c.ColName}};
{{- else if (eq $c.Operation 1)}}
-- 新增字段 {{$c.ColName}} 
ALTER TABLE {{$.Name}} add COLUMN {{$c.ColName}} {{$c|dbcolType}} {{$c|seq}}  {{$c|isNull}} {{$c|defaultValue}}  ;
{{$c|colComment}}

{{- else if (eq $c.Operation 2)}}
-- 修改字段 {{$c.ColName}} 
ALTER TABLE {{$.Name}} MODIFY {{$c.ColName}} {{$c|dbcolType}} {{$c|seq}} {{$c|isNull}}  {{$c|defaultValue}}  ;
{{$c|colComment}}

{{- end}}
{{- end}}


{{- range $i,$c:=.DiffIdxs}}
{{- if and (eq $c.Operation -1) ($c|isPk)}}
-- 删除主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP PRIMARY KEY;

{{- else if and (eq $c.Operation 1) ($c|isPk)}}
-- 新增主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} ADD {{$c|indexStr}};

{{- else if and (eq $c.Operation 2) ($c|isPk)}}
-- 修改主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP PRIMARY KEY;
ALTER TABLE {{$.Name}} ADD {{$c|indexStr}};

{{- else if and (eq $c.Operation -1) (or ($c|isIndex) ($c|isUNQ))}}
-- 删除索引 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP INDEX {{$c.Name}};
{{- else if and (eq $c.Operation 1) (or ($c|isIndex) ($c|isUNQ))}}
-- 新增索引 {{$c.Name}} 
ALTER TABLE {{$.Name}} ADD {{$c|indexStr}};
{{- else if and (eq $c.Operation 2) (or ($c|isIndex) ($c|isUNQ))}}
-- 修改索引 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP INDEX {{$c.Name}};
ALTER TABLE {{$.Name}} ADD {{$c|indexStr}};
{{- end}}
{{- end}}
`
