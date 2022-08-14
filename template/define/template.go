package define

import (
	"sync"

	"github.com/zhiyunliu/gluecli/model"
)

var (
	_tmpls = sync.Map{}
)

type Template interface {
	Name() string
	ReadPath(filePath string) (list *model.TmplTableList, err error)
	Translate(input interface{}) (string, error)
}

func Load(name string) Template {
	tmpv, ok := _tmpls.Load(name)
	if !ok {
		return nil
	}

	tmpl, ok := tmpv.(Template)
	if !ok {
		return nil
	}
	return tmpl
}

func Registry(tmpl Template) {
	_tmpls.Store(tmpl.Name(), tmpl)
}
