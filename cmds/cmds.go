package cmds

import (
	"github.com/urfave/cli"
)

var (
	calls = []RegisterCallback{}
)

type RegisterCallback func(app *cli.App)

func Bind(app *cli.App) {
	for _, call := range calls {
		call(app)
	}
}

func Register(call RegisterCallback) {
	calls = append(calls, call)
}
