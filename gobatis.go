package gobatis

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"
	"time"
)

var (
	ptnParam    = regexp.MustCompile(`#{(.*?)}`)
	ptnParamVar = regexp.MustCompile(`\${(.*?)}`)
)

type GoBatis struct {
	db      *sqlx.DB
	mappers map[string]*XmlParse
}

func NewGoBatis(cfg DbConfig) *GoBatis {
	o := &GoBatis{}

	o.db = sqlx.MustOpen(cfg.Driver, cfg.Dsn)
	o.db.SetMaxIdleConns(cfg.MaxIdleConnections)
	o.db.SetMaxOpenConns(cfg.MaxOpenConnections)
	o.db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleConnections) * time.Second)
	o.db.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)

	o.mappers = make(map[string]*XmlParse)

	return o
}

func (p *GoBatis) GetDb() *sqlx.DB {
	return p.db
}

func (p *GoBatis) LoadFromBytes(name string, bytes []byte) error {
	parser := &XmlParse{}
	if err := parser.LoadFromBytes(bytes); err != nil {
		return err
	}

	p.mappers[name] = parser
	return nil
}

func (p *GoBatis) LoadFromStr(name, str string) error {
	return p.LoadFromBytes(name, []byte(str))
}

func (p *GoBatis) LoadFromFile(name, filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return p.LoadFromBytes(name, b)
}

func (p *GoBatis) exec(selector string, args map[string]interface{}, fun func(stmt *sqlx.NamedStmt) error) error {
	s, err := p.getSelector(selector)
	if err != nil {
		return err
	}

	parser, ok := p.mappers[s.Name]
	if !ok {
		return errors.New("XML file \"" + s.Name + "\" is not exists!")
	}

	tsql, err := parser.Query(s.Id, args)
	if err != nil {
		return err
	}
	sugar.Info("raw SQL => ", tsql)

	stmt, err1 := p.db.PrepareNamed(p.bindVar(tsql, args))
	if err1 != nil {
		return err1
	}
	sugar.Info("prepared SQL => ", stmt.QueryString)

	return fun(stmt)
}

func (p *GoBatis) QueryObject(value interface{}, selector string, args map[string]interface{}) error {
	return p.exec(selector, args, func(stmt *sqlx.NamedStmt) error {
		return stmt.Get(value, args)
	})
}

func (p *GoBatis) QueryObjects(value interface{}, selector string, args map[string]interface{}) error {
	return p.exec(selector, args, func(stmt *sqlx.NamedStmt) error {
		return stmt.Select(value, args)
	})
}

func (p *GoBatis) update(selector string, args map[string]interface{}) (int64, error) {
	var n int64 = 0
	err := p.exec(selector, args, func(stmt *sqlx.NamedStmt) error {
		rs, err := stmt.Exec(args)
		if err != nil {
			return err
		}

		if strings.ToUpper(stmt.QueryString[0:6]) == "INSERT" {
			n, err = rs.LastInsertId()
		} else {
			n, err = rs.RowsAffected()
			sugar.Info(n)
		}
		return err
	})
	return n, err
}

func (p *GoBatis) Insert(selector string, args map[string]interface{}) (int64, error) {
	return p.update(selector, args)
}

func (p *GoBatis) Update(selector string, args map[string]interface{}) (int64, error) {
	return p.update(selector, args)
}

func (p *GoBatis) Delete(selector string, args map[string]interface{}) (int64, error) {
	return p.update(selector, args)
}

func (p *GoBatis) getSelector(selector string) (*selectorEntity, error) {
	arr := strings.Split(selector, ".")
	if len(arr) != 2 {
		return nil, errors.New("Selector \"" + selector + "\" is not exists!")
	}

	return &selectorEntity{
		Name: arr[0],
		Id:   arr[1],
	}, nil
}

func (p *GoBatis) bindVar(tsql string, args map[string]interface{}) string {
	tsql = ptnParam.ReplaceAllStringFunc(tsql, func(s string) string {
		return ":" + s[2:len(s)-1]
	})

	if args == nil {
		return tsql
	}

	return ptnParamVar.ReplaceAllStringFunc(tsql, func(s string) string {
		value := args[s[2:len(s)-1]]
		if value == nil || reflect.TypeOf(value).String() != "string" {
			return ""
		}
		return value.(string)
	})
}
