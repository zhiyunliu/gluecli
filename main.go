package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/zhiyunliu/gluecli/cmds"
	_ "github.com/zhiyunliu/gluecli/cmds/database"
)

func main() {
	var app = cli.NewApp()
	app.Version = "v0.0.1"
	cmds.Bind(app)
	app.Run(os.Args)
}
