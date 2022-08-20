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
go 

{{.|generateIdx}} 

go 

{{.|generateComment}}

go 
  `

const TmplDiffSQLModify = `
{{ $empty := "" -}}

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
{{- range $i,$dp:=$c.ColDiffPart}}
{{if (eq $dp 1)}}ALTER TABLE {{$.Name}} ALTER COLUMN {{$c.ColName}} {{$c|dbcolType}}  {{$c|isNull}};{{end}}
{{if (eq $dp 2)}}ALTER TABLE {{$.Name}} ADD CONSTRAINT  DF_{{$.Name}}_{{$c.ColName}}  DEFAULT ({{$c|alterDefaultValue}}) FOR {{$c.ColName}};{{end}}
{{if (eq $dp 3)}}{{$c|colComment}};{{end}}
{{- end}}
{{- end}}
{{- end}}


 
{{- range $i,$c:=.DiffIdxs}}
{{- if and (eq $c.Operation -1) ($c|isPk)}}
-- 删除主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP CONSTRAINT {{$c.Name}};
{{- else if and (eq $c.Operation 1) ($c|isPk)}}
-- 新增主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} ADD  CONSTRAINT $c.Name PRIMARY KEY CLUSTERED ({{$c|indexCols}});
{{- else if and (eq $c.Operation 2) ($c|isPk)}}
-- 修改主键 {{$c.Name}} 
ALTER TABLE {{$.Name}} DROP CONSTRAINT {{$c.Name}};
ALTER TABLE {{$.Name}} ADD  CONSTRAINT {{$c.Name}} PRIMARY KEY CLUSTERED ({{$c|indexCols}});
{{- else if and (eq $c.Operation -1)}}
-- 删除Index/UNQUE {{$c.Name}} 
DROP INDEX {{$c.Name}} ON  {{$.Name}};
{{- else if and (eq $c.Operation 1) ($c|isIDX)}}
-- 新增IDX索引 {{$c.Name}} 
CREATE NONCLUSTERED INDEX {{$c.Name}} ON {{$.Name}} ({{$c|indexCols}});
{{- else if and (eq $c.Operation 1) ($c|isUNQ)}}
-- 新增UNQ索引 {{$c.Name}} 
CREATE UNIQUE NONCLUSTERED INDEX {{$c.Name}} ON {{$.Name}} ({{$c|indexCols}});
{{- else if and (eq $c.Operation 2) ($c|isIDX)}}
-- 修改IDX索引 {{$c.Name}} 
DROP INDEX {{$c.Name}} ON  {{$.Name}};
CREATE NONCLUSTERED INDEX {{$c.Name}} ON {{$.Name}} ({{$c|indexCols}});
{{- else if and (eq $c.Operation 2) ($c|isUNQ)}}
-- 修改UNQ索引 {{$c.Name}} 
DROP INDEX {{$c.Name}} ON  {{$.Name}};
CREATE UNIQUE NONCLUSTERED INDEX {{$c.Name}} ON {{$.Name}} ({{$c|indexCols}});

{{- end}}
{{- end}}
`
