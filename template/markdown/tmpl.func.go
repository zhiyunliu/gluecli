package markdown

import (
	"github.com/micro-plat/lib4go/types"
)

type callHanlder func(string) string

func getfuncs(tp string) map[string]interface{} {
	return map[string]interface{}{
		"isTrue": types.GetBool,
	}
}
