package markdown

const TmplDictionary = `
###  {{.Comment}}[{{.Name}}]
| 字段名       | 类型       | 默认值  | 为空 |   约束    | 描述                                |
| ------------| -----------| :----: | :--: | :-------: | :---------------------------------|
{{range $i,$c:=.Cols -}}
| {{$c.ColName}} | {{$c|dbcolType}} |  {{$c.Default}} |  {{$c|isNull}} | {{$c|colCondition}} | {{$c.Comment}} |
{{end -}}

`
