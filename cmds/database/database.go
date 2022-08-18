package database

import (
	"github.com/urfave/cli"
	"github.com/zhiyunliu/gluecli/cmds"
)

func init() {
	cmds.Register(func(app *cli.App) {
		cmd := buildCommond()
		app.Commands = append(app.Commands, cmd)
	})
}

func buildCommond() (cmd cli.Command) {
	cmd = cli.Command{
		Name:  "db",
		Usage: "数据库结构文件",
	}
	subCmds := cli.Commands{
		buildSchemeCreateCmd(),
		buildSchemeDiffCmd(),
		//	buildSchemeDicCmd(),
	}
	cmd.Subcommands = subCmds
	return
}
