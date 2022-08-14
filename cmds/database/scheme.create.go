package database

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zhiyunliu/gluecli/database"
	"github.com/zhiyunliu/gluecli/logs"
	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/gluecli/template"
	"github.com/zhiyunliu/gluecli/xfile"
)

type schemeCreateOptions struct {
	DbType         string
	MdFilePath     string
	OutputPath     string
	TableName      string
	IncludeDrop    bool
	IncludeSeqFile bool
	NeedCoverFile  bool
}

func buildSchemeCreateCmd() cli.Command {
	opts := &schemeCreateOptions{}
	cmd := cli.Command{
		Name:  "create",
		Usage: "创建数据库结构文件",
		Action: func(ctx *cli.Context) (err error) {
			return createScheme(ctx, opts)
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "file,f", Destination: &opts.MdFilePath, Usage: `-md文件`},
			cli.StringFlag{Name: "out,o", Destination: &opts.OutputPath, Usage: `-输出路径`},
			cli.StringFlag{Name: "dbtype,db", Destination: &opts.DbType, Usage: `-数据库类型(mysql,mssql,oracle)`},
			cli.StringFlag{Name: "table,t", Destination: &opts.TableName, Usage: `-表名称`},
			cli.BoolFlag{Name: "drop,d", Destination: &opts.IncludeDrop, Usage: `-包含表删除语句`},
			cli.BoolFlag{Name: "seqfile,s", Destination: &opts.IncludeSeqFile, Usage: `-包含序列文件`},
			cli.BoolFlag{Name: "cover,v", Destination: &opts.NeedCoverFile, Usage: `-文件存在时覆盖`},
		},
	}
	return cmd
}

func createScheme(c *cli.Context, opts *schemeCreateOptions) (err error) {
	if len(c.Args()) == 0 {
		return fmt.Errorf("未指定markdown文件")
	}
	if len(c.Args()) < 2 {
		return fmt.Errorf("未指定输出路径")
	}

	var tblist *model.TmplTableList
	tblist, err = template.ReadPath(opts.MdFilePath)
	if err != nil {
		return err
	}

	//是否删除表
	tblist.DropTable(opts.IncludeDrop)
	//过滤数据表
	tblist.FilterTable(opts.TableName)
	//排除表(表名包含!的配置)
	tblist.Exclude()

	//循环创建表
	outpath := opts.OutputPath
	for _, tb := range tblist.Tables {
		//创建文件
		path := xfile.GetSchemePath(outpath, tb.Name)

		//转换文件
		content, err := database.BuildScheme(opts.DbType, tb)
		// content, err := tmpl.Translate(tmpl.SQLTmpl, opts.DbType, tb)
		if err != nil {
			return err
		}
		fs, err := xfile.Create(path, opts.NeedCoverFile)
		if err != nil {
			return err
		}
		logs.Info("生成文件:", path)
		if _, err := fs.Write([]byte(content)); err != nil {
			return err
		}
	}

	return nil
}
