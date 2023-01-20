package database

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zhiyunliu/gluecli/consts/enums/difffrom"
	"github.com/zhiyunliu/gluecli/database"
	"github.com/zhiyunliu/gluecli/logs"
	"github.com/zhiyunliu/gluecli/template"
	"github.com/zhiyunliu/gluecli/xfile"
)

type schemeDiffOptions struct {
	DbType        string
	From          string //file,server
	MdFilePathA   string
	MdFilePathB   string
	OutputPath    string
	TableName     string
	NeedCoverFile bool
}

func buildSchemeDiffCmd() cli.Command {
	opts := &schemeDiffOptions{}
	cmd := cli.Command{
		Name:  "diff",
		Usage: "创建数据库结构差异文件",
		Action: func(ctx *cli.Context) (err error) {
			return createDiff(ctx, opts)
		},
		Flags: []cli.Flag{
			cli.StringFlag{Name: "dbtype,db", Destination: &opts.DbType, Usage: `-数据库类型(mysql,mssql,oracle)`},
			cli.StringFlag{Name: "from", Destination: &opts.From, Value: difffrom.File, Usage: `-数据库类型(mysql,mssql,oracle)`},
			cli.StringFlag{Name: "filesrc,fa", Destination: &opts.MdFilePathA, Usage: `-src文件`},
			cli.StringFlag{Name: "filedest,fb", Destination: &opts.MdFilePathB, Usage: `-dest文件`},
			cli.StringFlag{Name: "out,o", Destination: &opts.OutputPath, Usage: `-输出路径`},
			cli.StringFlag{Name: "table,t", Destination: &opts.TableName, Usage: `-表名称`},
			cli.BoolFlag{Name: "cover,v", Destination: &opts.NeedCoverFile, Usage: `-文件存在时覆盖`},
		},
	}
	return cmd
}

//createDiff 生成数据库结构
func createDiff(c *cli.Context, opts *schemeDiffOptions) (err error) {
	if opts.MdFilePathA == "" || opts.MdFilePathB == "" {
		return fmt.Errorf("未指定需要对比的源. --fa,--fb")
	}

	targetTbs, err := template.ReadPath(opts.MdFilePathA)
	if err != nil {
		return err
	}
	sourceTbs, err := template.ReadPath(opts.MdFilePathB)
	if err != nil {
		return err
	}
	//过滤表
	sourceTbs.FilterTable(opts.TableName)
	sourceTbs.Exclude()

	//过滤表
	targetTbs.FilterTable(opts.TableName)
	targetTbs.Exclude()

	diff := sourceTbs.Diff(targetTbs)

	if len(diff.Tables) == 0 {
		return fmt.Errorf("两个文件是一样的")
	}

	for _, tb := range diff.Tables {
		//创建文件
		path := xfile.GetSchemePath(opts.OutputPath, tb.Name)

		//翻译文件
		content, err := database.Diff(opts.DbType, tb)
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

func getDiffTableData() {

}
