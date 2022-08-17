package main

import (
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/zhiyunliu/gluecli/cmds"
	_ "github.com/zhiyunliu/gluecli/cmds/database"
	"github.com/zhiyunliu/golibs/xlog/log"
)

func main() {
	var app = cli.NewApp()
	app.Version = "v0.0.1"
	app.ExitErrHandler = func(context *cli.Context, err error) {
		log.Error(err)
		time.Sleep(time.Second * 2)
	}
	cmds.Bind(app)

	app.Run(os.Args)
}
