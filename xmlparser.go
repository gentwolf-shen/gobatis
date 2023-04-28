package gobatis

import (
	"errors"
	"strings"

	"github.com/beevik/etree"
)

type XmlParser struct {
	doc *etree.Document
}

func (xml *XmlParser) LoadFromBytes(bytes []byte) error {
	xml.doc = etree.NewDocument()
	return xml.doc.ReadFromBytes(bytes)
}

// 从XML中解析出SQL，及绑定的参数
func (xml *XmlParser) Query(selector *selectorEntity, inputValue map[string]interface{}) (*Queryer, error) {
	q := &Queryer{}

	item := xml.doc.FindElement("./mapper/*[@id='" + selector.Id + "']")
	if item == nil {
		return q, errors.New("XML id \"" + selector.Name + "." + selector.Id + "\" is not exists")
	}

	q.Value = make(map[string]interface{})
	parser := NewXmlParserBuild(inputValue, q.Value)
	q.StatementType = strings.ToUpper(item.SelectAttrValue("statementType", "PREPARED"))

	var err error
	q.Sql, err = parser.Build(item)
	return q, err
}
