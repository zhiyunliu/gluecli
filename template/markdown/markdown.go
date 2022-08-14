package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/zhiyunliu/gluecli/model"
	"github.com/zhiyunliu/gluecli/template/define"
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
	totalTableList := &model.TmplTableList{
		Map: make(map[string]bool),
	}
	for _, fn := range fns {
		newTable, err := readFile(fn)
		if err != nil {
			return nil, err
		}
		for key := range newTable.Map {
			if _, ok := totalTableList.Map[key]; ok {
				return nil, fmt.Errorf("存在相同的表名：%s", key)
			}
			totalTableList.Map[key] = true
		}
		totalTableList.Tables = append(totalTableList.Tables, newTable.Tables...)
	}
	return totalTableList, nil
}

func (m *markdown) Translate(input interface{}) (string, error) {
	var tmpl = template.New("table").Funcs(getfuncs(tp))
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
