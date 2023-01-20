package database

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/zhiyunliu/gluecli/consts/enums/dbtype"
	"github.com/zhiyunliu/gluecli/consts/enums/tmpltype"
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
	opts := &schemeDicOptions{
		DbType: "",
	}
	cmd := cli.Command{
		Name:  "file",
		Usage: "根据数据库生成文件",
		Action: func(ctx *cli.Context) (err error) {
			return dicScheme(ctx, opts)
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "dbconn,dc", Destination: &opts.DbConn, Usage: `-数据库连接`},
			cli.StringFlag{Name: "out,o", Destination: &opts.OutputPath, Usage: `-输出路径`},
			cli.StringFlag{Name: "dbtype,db", Destination: &opts.DbType, Value: dbtype.MsSql, Usage: `-数据库类型(mysql,mssql,oracle)`},
			cli.StringFlag{Name: "table,t", Destination: &opts.TableName, Usage: `-表名称`},
			cli.StringFlag{Name: "tmpl,tpl", Destination: &opts.TmplType, Value: tmpltype.Markdown, Usage: `-数据文件类型(md/doc/docx)`},
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

	//生成文件
	path := filepath.Join(opts.OutputPath, fmt.Sprintf("./%s.%s", opts.DbType, opts.TmplType))
	logs.Info("生成文件:", path, opts.NeedCoverFile)
	fs, err := xfile.Create(path, opts.NeedCoverFile)
	if err != nil {
		return err
	}
	defer fs.Close()

	for _, tbl := range dbInfo.Tables {
		//翻译文件
		err := template.Translate(fs, opts.TmplType, opts.DbType, tbl)
		if err != nil {
			return err
		}
	}

	return nil
}
