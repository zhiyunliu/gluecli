package template

import (
	"fmt"
	"path/filepath"

	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/gluecli/template/define"
	_ "github.com/zhiyunliu/gluecli/template/markdown"
)

func ReadPath(filePath string) (tbList *model.TmplTableList, err error) {
	if filePath == "" {
		err = fmt.Errorf("filePath 不能为空")
		return
	}
	fileType := getProviderName(filePath)
	if fileType == "" {
		err = fmt.Errorf("filePath 格式错误,需要以***.md/***.doc/***.docx结尾")
		return
	}
	tmpl := define.Load(fileType)
	if tmpl == nil {
		err = fmt.Errorf("未注册类型:%s", fileType)
		return
	}
	return tmpl.ReadPath(filePath)
}

func Translate(tmpl string, input interface{}) (content string, err error) {
	impl := define.Load(tmpl)
	if impl == nil {
		err = fmt.Errorf("未注册类型:%s", tmpl)
		return
	}
	return impl.Translate(input)
}

func getProviderName(filePath string) string {
	return filepath.Ext(filePath)
}
