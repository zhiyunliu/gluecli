package markdown

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/gluecli/template/define"
)

const (
	Template = "markdown"
)

func init() {
	define.Registry(&markdown{})
}

type markdown struct{}

func (m *markdown) Name() string {
	return Template
}
func (m *markdown) ReadPath(filePath string) (list *model.TmplTableList, err error) {
	fns := getMatchFiles(filePath)
	//读取文件
	totalTableList := model.NewTableList()
	for _, fn := range fns {
		newTable, err := readFile(fn)
		if err != nil {
			return nil, err
		}

		for i := range newTable.Tables {
			err = totalTableList.Append(newTable.Tables[i])
			if err != nil {
				return nil, err
			}
		}
	}
	return totalTableList, nil
}

func (m *markdown) Translate(dbType string, input interface{}) (string, error) {
	var tmpl = template.New("table").Funcs(getfuncs(dbType))
	np, err := tmpl.Parse(TmplDictionary)
	if err != nil {
		return "", err
	}
	buff := bytes.NewBufferString("")
	if err := np.Execute(buff, input); err != nil {
		return "", err
	}
	return strings.Replace(strings.Replace(buff.String(), "{###}", "`", -1), "&#39;", "'", -1), nil
}
