package gobatis

import (
	"errors"
	"github.com/beevik/etree"
	"reflect"
	"regexp"
	"strings"
)

var (
	ptnOp   = regexp.MustCompile(`[\s]*(.*?)[\s]*(!=|==|>=|<=|>|<|=)[\s]*(.*)[\s]*`)
	ptnTest = regexp.MustCompile(`[\s]+(AND|and|OR|or)[\s]+`)
)

type XmlParse struct {
	doc *etree.Document
}

func (p *XmlParse) LoadFromBytes(bytes []byte) error {
	p.doc = etree.NewDocument()
	return p.doc.ReadFromBytes(bytes)
}

func (p *XmlParse) Query(id string, params map[string]interface{}) (string, error) {
	item := p.doc.FindElement("./mapper/*[@id='" + id + "']")
	if item == nil {
		return "", errors.New("XML id \"" + id + "\" is not exists")
	}

	var builder strings.Builder
	for _, child := range item.Child {
		tmp := ""
		switch reflect.TypeOf(child).String() {
		case "*etree.CharData":
			tmp = strings.TrimSpace(child.(*etree.CharData).Data)
		case "*etree.Element":
			el := child.(*etree.Element)
			switch el.Tag {
			case "where":
				tmp = p.buildWhere(el, params)
			case "set":
				tmp = p.buildSet(el, params)
			case "if":
				tmp = p.buildIf(el, params)
			}
		}
		if tmp != "" {
			builder.WriteString(tmp)
			builder.WriteString(" ")
		}
	}

	return builder.String(), nil
}

func (p *XmlParse) buildWhere(item *etree.Element, params map[string]interface{}) string {
	wheres := p.buildWhereOrSet(item, params)
	if len(wheres) > 0 {
		return "WHERE " + strings.Join(wheres, " AND ")
	}
	return ""
}

func (p *XmlParse) buildSet(item *etree.Element, params map[string]interface{}) string {
	sets := p.buildWhereOrSet(item, params)
	if len(sets) > 0 {
		return "SET " + strings.Join(sets, ", ")
	}
	return ""
}

func (p *XmlParse) buildWhereOrSet(item *etree.Element, params map[string]interface{}) []string {
	ifElements := item.FindElements("if")
	wheres := make([]string, len(ifElements))
	index := 0
	for _, el := range ifElements {
		str := strings.TrimSpace(el.SelectAttrValue("test", ""))
		if p.parseTest(str, params) {
			wheres[index] = p.filter(el.Text())
			index++
		}
	}
	return wheres[0:index]
}

func (p *XmlParse) buildIf(item *etree.Element, params map[string]interface{}) string {
	str := strings.TrimSpace(item.SelectAttrValue("test", ""))
	if p.parseTest(str, params) {
		return p.filter(item.Text())
	}
	return ""
}

func (p *XmlParse) filter(str string) string {
	str = strings.Replace(str, "&lt;", "<", -1)
	str = strings.Replace(str, "&gt;", ">", -1)
	return strings.TrimSpace(str)
}

func (p *XmlParse) parseTest(str string, args map[string]interface{}) bool {
	if str == "" {
		return true
	}

	str = ptnTest.ReplaceAllStringFunc(str, strings.ToUpper)
	bl := false

	if strings.Contains(str, " OR ") {
		for _, item := range strings.Split(str, " OR ") {
			if bl = p.testVal(item, args); bl {
				break
			}
		}
		return bl
	}

	for _, item := range strings.Split(str, " AND ") {
		if bl = p.testVal(item, args); !bl {
			break
		}
	}

	return bl
}

func (p *XmlParse) testVal(testStr string, args map[string]interface{}) bool {
	key, op, testValue := p.getTestSegment(testStr)
	value, ok := args[key]
	if !ok {
		return false
	}

	return compare(value, op, testValue)
}

func (p *XmlParse) getTestSegment(testStr string) (string, string, string) {
	segments := ptnOp.FindStringSubmatch(testStr)
	if len(segments) != 4 {
		return "", "", ""
	}

	return segments[1], segments[2], segments[3]
}
