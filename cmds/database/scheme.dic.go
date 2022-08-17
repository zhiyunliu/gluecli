package database

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/zhiyunliu/gluecli/database"
	"github.com/zhiyunliu/gluecli/logs"
	"github.com/zhiyunliu/gluecli/template"
	"github.com/zhiyunliu/gluecli/xfile"
)

type schemeDicOptions struct {
	DbType        string
	DbConn        string
	OutputPath    string
	TableName     string
	TmplType      string
	NeedCoverFile bool
}

func buildSchemeDicCmd() cli.Command {
	opts := &schemeDicOptions{}
	cmd := cli.Command{
		Name:  "file",
		Usage: "根据数据库生成文件",
		Action: func(ctx *cli.Context) (err error) {
			return dicScheme(ctx, opts)
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "dbconn,dc", Destination: &opts.DbConn, Usage: `-md文件`},
			cli.StringFlag{Name: "out,o", Destination: &opts.OutputPath, Usage: `-输出路径`},
			cli.StringFlag{Name: "dbtype,dt", Destination: &opts.DbType, Usage: `-数据库类型(mysql,mssql,oracle)`},
			cli.StringFlag{Name: "table,t", Destination: &opts.TableName, Usage: `-表名称`},
			cli.StringFlag{Name: "tmpl,tpl", Destination: &opts.TmplType, Usage: `-数据文件类型(md/doc/docx)`},
			cli.BoolFlag{Name: "cover,v", Destination: &opts.NeedCoverFile, Usage: `-文件存在时覆盖`},
		},
	}
	return cmd
}

func dicScheme(c *cli.Context, opts *schemeDicOptions) (err error) {

	dbInfo, err := database.GetDbInfo(opts.DbType, opts.DbConn)
	if err != nil {
		return
	}

	//循环创建表
	content := ""
	for _, tb := range dbInfo.Tables {
		//翻译文件
		ct, err := template.Translate(opts.TmplType, opts.DbType, tb)
		if err != nil {
			return err
		}
		content += ct
	}

	//生成文件
	path := filepath.Join(opts.OutputPath, fmt.Sprintf("./%s.md", opts.DbType))
	fs, err := xfile.Create(path, opts.NeedCoverFile)
	if err != nil {
		return err
	}
	logs.Info("生成文件:", path)
	fs.WriteString(content)
	fs.Close()

	return nil
}

const MdDictionaryTPL = `
{{ $empty := "" -}}
###  {{.Desc}}[{{.Name}}]

| 字段名       | 类型       | 默认值  | 为空 |   约束    | 描述                                |
| ------------| -----------| :----: | :--: | :-------: | :---------------------------------|
{{range $i,$c:=.Rows -}}
| {{$c.Name}} | {{$c.Type}}{{if ne $c.LenStr $empty}}({{$c.LenStr}}){{end}}|{{$c.Def}}|{{$c|isMDNull}}| {{$c.Con}} | {{$c.Desc}}|
{{end -}}
`
