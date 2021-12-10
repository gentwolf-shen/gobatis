package gobatis

import (
	"errors"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/beevik/etree"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type XmlParseForeach struct {
	separatorAttr  string
	openAttr       string
	closeAttr      string
	indexAttr      string
	itemAttr       string
	collectionAttr string
	index          int
	ptnParamIndex  *regexp.Regexp
	ptnParamItem   *regexp.Regexp
	convertedValue map[string]interface{}
	outputValue    map[string]interface{}
}

// 构造SQL
func (b *XmlParseForeach) Build(el *etree.Element, inputValue map[string]interface{}, outputValue map[string]interface{}) (string, error) {
	b.getAttr(el)
	if err := b.convertValue(inputValue); err != nil {
		return "", err
	}

	b.outputValue = outputValue

	b.ptnParamIndex = regexp.MustCompile(`[$#]{` + b.indexAttr + `}`)
	b.ptnParamItem = regexp.MustCompile(`[$#]{` + b.itemAttr + `}`)

	index := 0
	arr := make([]string, len(b.convertedValue))

	text := strings.TrimSpace(el.Text())
	for k, v := range b.convertedValue {
		str := b.replaceIndex(text, k, v)
		str = b.replaceItem(str, v)
		arr[index] = str
		index++
	}

	if index == 0 {
		return "", nil
	}

	return b.openAttr + strings.Join(arr, b.separatorAttr) + b.closeAttr, nil
}

// 获取foreach上的属性
func (b *XmlParseForeach) getAttr(el *etree.Element) {
	b.separatorAttr = el.SelectAttrValue("separator", "")
	b.openAttr = el.SelectAttrValue("open", "")
	b.closeAttr = el.SelectAttrValue("close", "")
	b.indexAttr = el.SelectAttrValue("index", "")
	b.itemAttr = el.SelectAttrValue("item", "")
	b.collectionAttr = el.SelectAttrValue("collection", "")
}

// 将collection关联的值转换为map[string]interface{}格式，方便处理
// 如果是slice，转换后顺序可能会不对
func (b *XmlParseForeach) convertValue(inputValue interface{}) error {
	if b.collectionAttr == "" {
		return errors.New("the attribute \"collection\" is not exists")
	}

	p, err := expr.Eval(b.collectionAttr, inputValue)
	if err != nil {
		return err
	}

	if p == nil {
		logger.Warn("foreach \"", b.collectionAttr, "\" value is nil")
		return nil
	}

	b.convertedValue = make(map[string]interface{})
	values := reflect.ValueOf(p)
	kind := reflect.TypeOf(p).Kind()
	if kind == reflect.Map {
		for _, v := range values.MapKeys() {
			b.convertedValue[v.String()] = values.MapIndex(v).Interface()
		}
	} else if kind == reflect.Slice {
		for i := 0; i < values.Len(); i++ {
			b.convertedValue[strconv.Itoa(i)] = values.Index(i).Interface()
		}
	}

	return nil
}

// 替换index的标记
func (b *XmlParseForeach) replaceIndex(str string, key string, value interface{}) string {
	if b.indexAttr == "" {
		return str
	}

	return b.ptnParamIndex.ReplaceAllStringFunc(str, func(s string) string {
		if s[0:1] == "$" {
			return key
		}

		str := b.collectionAttr + "_" + key
		b.outputValue[str] = value

		return str
	})
}

// 替换item的标记
func (b *XmlParseForeach) replaceItem(str string, value interface{}) string {
	if b.itemAttr == "" {
		return str
	}

	var p1 *vm.Program
	return b.ptnParamItem.ReplaceAllStringFunc(str, func(s string) string {
		input := s[2 : len(s)-1]
		if p1 == nil {
			var err error
			if p1, err = expr.Compile(input); err != nil {
				logger.Error("eval ", input, " error: ", err)
			}
		}

		output, err := expr.Run(p1, map[string]interface{}{"item": value})
		if err != nil {
			logger.Error("run ", input, " error: ", err)
		}

		if s[0:1] == "$" {
			return toStr(output)
		}

		tmp := fmt.Sprintf("%s_%d", b.collectionAttr, b.index)
		b.index++
		b.outputValue[tmp] = output

		return ":" + tmp
	})
}
