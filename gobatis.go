package gobatis

import (
	"errors"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"strings"
	"time"
)

type GoBatis struct {
	db      *sqlx.DB
	mappers map[string]*XmlParser
}

func NewGoBatis(cfg DbConfig) *GoBatis {
	o := &GoBatis{}

	o.db = sqlx.MustOpen(cfg.Driver, cfg.Dsn)
	o.db.SetMaxIdleConns(cfg.MaxIdleConnections)
	o.db.SetMaxOpenConns(cfg.MaxOpenConnections)
	o.db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleConnections) * time.Second)
	o.db.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)

	o.mappers = make(map[string]*XmlParser)

	return o
}

func (p *GoBatis) GetDb() *sqlx.DB {
	return p.db
}

func (p *GoBatis) LoadFromBytes(name string, bytes []byte) error {
	parser := &XmlParser{}
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

func (p *GoBatis) exec(selector string, inputValue map[string]interface{}, fun func(stmt *sqlx.NamedStmt, outputValue map[string]interface{}) error) error {
	s, err := p.getSelector(selector)
	if err != nil {
		return err
	}

	parser, ok := p.mappers[s.Name]
	if !ok {
		return errors.New("XML file \"" + s.Name + "\" is not exists!")
	}

	tsql, outputValue, err := parser.Query(s.Id, inputValue)
	if err != nil {
		return err
	}

	logger.Debug("raw SQL => \n", tsql)

	stmt, err1 := p.db.PrepareNamed(tsql)
	if err1 != nil {
		return err1
	}
	logger.Debug("prepared SQL => \n", stmt.QueryString)

	return fun(stmt, outputValue)
}

func (p *GoBatis) QueryObject(value interface{}, selector string, inputValue map[string]interface{}) error {
	return p.exec(selector, inputValue, func(stmt *sqlx.NamedStmt, outputValue map[string]interface{}) error {
		return stmt.Get(value, outputValue)
	})
}

func (p *GoBatis) QueryObjects(value interface{}, selector string, inputValue map[string]interface{}) error {
	return p.exec(selector, inputValue, func(stmt *sqlx.NamedStmt, outputValue map[string]interface{}) error {
		return stmt.Select(value, outputValue)
	})
}

func (p *GoBatis) update(selector string, inputValue map[string]interface{}) (int64, error) {
	var n int64 = 0
	err := p.exec(selector, inputValue, func(stmt *sqlx.NamedStmt, outputValue map[string]interface{}) error {
		rs, err := stmt.Exec(outputValue)
		if err != nil {
			return err
		}

		if strings.ToUpper(stmt.QueryString[0:6]) == "INSERT" {
			n, err = rs.LastInsertId()
		} else {
			n, err = rs.RowsAffected()
		}
		return err
	})
	return n, err
}

func (p *GoBatis) Insert(selector string, inputValue map[string]interface{}) (int64, error) {
	return p.update(selector, inputValue)
}

func (p *GoBatis) Update(selector string, inputValue map[string]interface{}) (int64, error) {
	return p.update(selector, inputValue)
}

func (p *GoBatis) Delete(selector string, inputValue map[string]interface{}) (int64, error) {
	return p.update(selector, inputValue)
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
