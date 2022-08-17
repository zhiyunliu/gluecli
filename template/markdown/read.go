package markdown

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/micro-plat/lib4go/types"
	"github.com/zhiyunliu/gluecli/model"
)

//Line 每一行信息
type Line struct {
	Text    string
	LineNum int
}

//TableLine 表的每一行
type TableLine struct {
	Lines [][]*Line
}

//ReadFile 读取markdown文件并转换为MarkDownDB对象
func readFile(fileName string) (*model.TmplTableList, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bnr := bufio.NewReader(file)
	lines := loadLinesFromReader(bnr)
	return tableLine2Table(line2TableLine(lines))
}

func loadLinesFromReader(bnr *bufio.Reader) []*Line {
	lines := make([]*Line, 0, 64)
	num := 0
	for {
		num++
		line, err := bnr.ReadString('\n')
		if line == "" && (err != nil || io.EOF == err) {
			break
		}
		line = strings.Trim(strings.Trim(line, "\n"), "\t")
		if strings.TrimSpace(line) == "" {
			continue
		}
		lines = append(lines, &Line{Text: line, LineNum: num})
	}

	return lines
}

//lines2TableLine 数据行转变为以表为单个整体的数据行
func line2TableLine(lines []*Line) (tl TableLine) {
	dlines := []int{}
	for i, line := range lines {
		text := strings.TrimSpace(strings.Replace(line.Text, " ", "", -1))
		if text == "|字段名|类型|默认值|为空|约束|描述|" {
			dlines = append(dlines, i-1)
		}
		if len(dlines)%2 == 1 && strings.Count(text, "|") != 7 {
			dlines = append(dlines, i-1)
		}
	}
	if len(dlines)%2 == 1 {
		dlines = append(dlines, len(lines)-1)
	}
	//划分为以一张表为一个整体
	for i := 0; i < len(dlines); i = i + 2 {
		tl.Lines = append(tl.Lines, lines[dlines[i]:dlines[i+1]+1])
	}
	return tl
}

//tableLine2Table 表数据行变为表
func tableLine2Table(lines TableLine) (tables *model.TmplTableList, err error) {
	tables = &model.TmplTableList{Tables: make([]*model.TmplTable, 0, 1), Map: make(map[string]bool)}
	for _, tline := range lines.Lines {
		//markdown表格的表名，标题，标题数据区分行，共三行
		if len(tline) <= 3 {
			continue
		}
		var tb *model.TmplTable
		for i, line := range tline {
			if i == 0 {
				//获取表名，描述名称
				name, err := getTableName(line)
				if err != nil {
					return nil, err
				}
				tb = model.NewTmplTable(name, getTableDesc(line), getTableExtInfo(line))
				//tb.Name = getDbObjectName(line)
				continue
			}
			if i < 3 {
				continue
			}
			c, err := line2TableCol(line)
			if err != nil {
				return nil, err
			}

			if err := tb.AddCol(c); err != nil {
				return nil, err
			}
		}
		if tb != nil {
			if _, ok := tables.Map[tb.Name]; ok {
				return nil, fmt.Errorf("存在相同的表名：%s", tb.Name)
			}
			tables.Map[tb.Name] = true
			tables.Tables = append(tables.Tables, tb)
		}
	}
	return tables, nil
}

func line2TableCol(line *Line) (*model.TmplCol, error) {
	if strings.Count(line.Text, "|") != 7 {
		return nil, fmt.Errorf("表结构有误(行:%d)", line.LineNum)
	}
	colums := strings.Split(strings.Trim(line.Text, "|"), "|")
	if colums[0] == "" {
		return nil, fmt.Errorf("字段名称不能为空 %s(行:%d)", line.Text, line.LineNum)
	}

	coltp, len, decimalLen, err := getColType(line)
	if err != nil {
		return nil, err
	}
	c := &model.TmplCol{
		LineNum:    line.LineNum,
		ColName:    strings.TrimSpace(strings.Replace(colums[0], "&#124;", "|", -1)),
		ColType:    coltp,
		ColLen:     len,
		DecimalLen: decimalLen,
		Default:    strings.TrimSpace(strings.Replace(colums[2], "&#124;", "|", -1)),
		IsNull:     strings.TrimSpace(colums[3]),
		Condition:  strings.TrimSpace(colums[4]), // strings.Replace(strings.TrimSpace(colums[4]), " ", "", -1),
		Comment:    strings.TrimSpace(strings.Replace(colums[5], "&#124;", "|", -1)),
	}
	return c, nil
}

func getTableDesc(line *Line) string {
	reg := regexp.MustCompile(`[^\d\.|\s]?[^\x00-\xff]+[^\[]+`)
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return ""
	}
	return strings.TrimSpace(names[0])
}

func getTableExtInfo(line *Line) string {
	reg := regexp.MustCompile(`{(.*?)}$`)
	names := reg.FindStringSubmatch(line.Text)
	if len(names) == 0 {
		return ""
	}
	return strings.TrimSpace(names[1])
}

func getTableName(line *Line) (string, error) {
	if !strings.HasPrefix(line.Text, "###") {
		return "", fmt.Errorf("%d行表名称标注不正确，请以###开头:(%s)", line.LineNum, line.Text)
	}

	reg := regexp.MustCompile(`\[[\^]?[\w]+[,]?[\p{Han}A-Za-z0-9_]+\]`)
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return "", fmt.Errorf("未设置表名称或者格式不正确:%s(行:%d)，格式：### 描述[表名,菜单名]，菜单名可选", line.Text, line.LineNum)
	}
	s := strings.Split(strings.TrimRight(strings.TrimLeft(names[0], "["), "]"), ",")
	return s[0], nil
}

func getDbObjectName(line *Line) string {
	reg := regexp.MustCompile(`\@[\w]+`) //数据库对象名字
	names := reg.FindAllString(line.Text, -1)
	if len(names) == 0 {
		return ""
	}
	return fmt.Sprintf(`"%s"`, strings.TrimPrefix(names[0], "@"))
}

//类型，长度，小数长度，错误
func getColType(line *Line) (string, int, int, error) {
	colums := strings.Split(strings.Trim(line.Text, "|"), "|")
	if colums[0] == "" {
		return "", 0, 0, fmt.Errorf("字段名称不能为空 %s(行:%d)", line.Text, line.LineNum)
	}

	t := strings.TrimSpace(colums[1])
	reg := regexp.MustCompile(`[\w]+`)
	names := reg.FindAllString(t, -1)
	if len(names) == 0 || len(names) > 4 {
		return "", 0, 0, fmt.Errorf("未设置字段类型:%v(行:%d)", names, line.LineNum)
	}
	if len(names) == 1 {
		return t, 0, 0, nil
	}
	if len(names) == 2 {
		return t, types.GetInt(names[1]), 0, nil
	}
	return t, types.GetInt(names[1]), types.GetInt(names[2]), nil
}

func getMatchFiles(path string) (paths []string) {
	//路径是的具体文件
	_, err := os.Stat(path)
	if err == nil {
		return []string{path}
	}
	//查找匹配的文件
	dir, f := filepath.Split(path)

	regexName := fmt.Sprintf("^%s$", strings.Replace(strings.Replace(f, ".md", "\\.md", -1), "*", "(.+)?", -1))
	reg := regexp.MustCompile(regexName)

	if dir == "" {
		dir = "./"
	}
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		fname := f.Name()
		if strings.HasPrefix(fname, ".") || f.IsDir() {
			continue
		}
		if reg.Match([]byte(fname)) {
			paths = append(paths, filepath.Join(dir, fname))
		}
	}
	return paths
}
