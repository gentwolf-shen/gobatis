package gobatis

import (
	"github.com/antonmedv/expr"
	"github.com/beevik/etree"
	"reflect"
	"regexp"
	"strings"
)

var (
	ptnOutputVar     = regexp.MustCompile(`#{(.*?)}`)
	ptnOutputReplace = regexp.MustCompile(`\${(.*?)}`)
)

type XmlParserBuild struct {
	inputValue  map[string]interface{}
	outputValue map[string]interface{}
}

func NewXmlParserBuild(inputValue, outputValue map[string]interface{}) *XmlParserBuild {
	return &XmlParserBuild{
		inputValue:  inputValue,
		outputValue: outputValue,
	}
}

func (xml *XmlParserBuild) Build(element *etree.Element) (string, error) {
	var err error

	var builder strings.Builder
	for _, child := range element.Child {
		tmp := ""
		switch reflect.TypeOf(child).String() {
		case "*etree.CharData":
			tmp = xml.filter(strings.TrimSpace(child.(*etree.CharData).Data))
		case "*etree.Element":
			el := child.(*etree.Element)
			switch el.Tag {
			case "trim":
				if tmp, err = xml.buildTrim(el); err != nil {
					break
				}
			case "where":
				tmp = xml.buildWhere(el)
			case "set":
				tmp = xml.buildSet(el)
			case "if":
				tmp = xml.buildIf(el)
			case "foreach":
				b := XmlParseForeach{}
				tmp, err = b.Build(el, xml.inputValue, xml.outputValue)
				if err != nil {
					return "", err
				}
			}
		}
		if tmp != "" {
			builder.WriteString(tmp)
			builder.WriteString(" ")
		}
	}

	return strings.TrimSpace(builder.String()), err
}

func (xml *XmlParserBuild) buildTrim(el *etree.Element) (string, error) {
	tmp, err := xml.Build(el)
	if err != nil {
		return "", err
	}

	if tmp == "" {
		return tmp, err
	}

	if str := el.SelectAttrValue("prefixOverrides", ""); str != "" {
		for _, tag := range strings.Split(str, "|") {
			tmp = strings.TrimLeft(tmp, strings.TrimSpace(tag))
		}
	}

	if str := el.SelectAttrValue("suffixOverrides", ""); str != "" {
		for _, tag := range strings.Split(str, "|") {
			tmp = strings.TrimRight(tmp, strings.TrimSpace(tag))
		}
	}

	if tmp != "" {
		tmp = el.SelectAttrValue("prefix", "") + tmp + el.SelectAttrValue("suffix", "")
	}

	return tmp, nil
}

func (xml *XmlParserBuild) buildWhere(el *etree.Element) string {
	wheres := xml.buildWhereOrSet(el)
	if len(wheres) > 0 {
		return "WHERE " + strings.Join(wheres, " AND ")
	}
	return ""
}

func (xml *XmlParserBuild) buildSet(el *etree.Element) string {
	sets := xml.buildWhereOrSet(el)
	if len(sets) > 0 {
		return "SET " + strings.Join(sets, ", ")
	}
	return ""
}

func (xml *XmlParserBuild) buildWhereOrSet(el *etree.Element) []string {
	ifElements := el.FindElements("if")
	wheres := make([]string, len(ifElements))
	index := 0
	for _, el := range ifElements {
		str := strings.TrimSpace(el.SelectAttrValue("test", ""))
		if xml.parseTest(str, xml.inputValue) {
			wheres[index] = xml.filter(el.Text())
			index++
		}
	}
	return wheres[0:index]
}

func (xml *XmlParserBuild) buildIf(el *etree.Element) string {
	str := strings.TrimSpace(el.SelectAttrValue("test", ""))
	if xml.parseTest(str, xml.inputValue) {
		return xml.filter(el.Text())
	}
	return ""
}

func (xml *XmlParserBuild) filter(str string) string {
	str = strings.Replace(str, "&lt;", "<", -1)
	str = strings.Replace(str, "&gt;", ">", -1)

	str = ptnOutputVar.ReplaceAllStringFunc(str, func(s string) string {
		key := strings.TrimSpace(s[2 : len(s)-1])
		xml.outputValue[key] = xml.inputValue[key]
		return ":" + key
	})

	str = ptnOutputReplace.ReplaceAllStringFunc(str, func(s string) string {
		key := strings.TrimSpace(s[2 : len(s)-1])
		return toStr(xml.inputValue[key])
	})

	return strings.TrimSpace(str)
}

func (xml *XmlParserBuild) parseTest(str string, args map[string]interface{}) bool {
	if str == "" {
		return true
	}

	p, err := expr.Eval(str, args)
	if err != nil {
		logger.Error("eval ", str, " error: ", err)
		return false
	}

	valueType := reflect.ValueOf(p).Kind().String()
	if valueType == "bool" {
		return p.(bool)
	}

	logger.Warn("Bool is required by  \"", str+"\", but it is ", reflect.TypeOf(p))

	return false
}
