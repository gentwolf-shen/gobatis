package gobatis

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type GoBatis struct {
	lock    *sync.RWMutex
	db      *sqlx.DB
	mappers map[string]*XmlParser
	stmts   sync.Map
}

func NewGoBatis(cfg DbConfig) *GoBatis {
	o := &GoBatis{}

	o.db = sqlx.MustOpen(cfg.Driver, cfg.Dsn).Unsafe()
	o.db.SetMaxIdleConns(cfg.MaxIdleConnections)
	o.db.SetMaxOpenConns(cfg.MaxOpenConnections)
	o.db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleConnections) * time.Second)
	o.db.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)
	o.stmts = sync.Map{}
	o.lock = &sync.RWMutex{}
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

func (p *GoBatis) exec(selector string, inputValue map[string]interface{}, fun func(stmt *sqlx.NamedStmt, queryer *Queryer) error) error {
	s, err := p.getSelector(selector)
	if err != nil {
		return err
	}

	parser, ok := p.mappers[s.Name]
	if !ok {
		return errors.New("XML file \"" + s.Name + "\" is not exists!")
	}

	queryer, err := parser.Query(s.Id, inputValue)
	if err != nil {
		return err
	}

	logger.Debug("raw SQL => \n", queryer.Sql)

	if !strings.Contains(queryer.Sql, ":") {
		return fun(nil, queryer)
	}

	stmt, err := p.getStmt(queryer)
	if err != nil {
		return err
	}

	return fun(stmt, queryer)
}

func (p *GoBatis) QueryObject(dest interface{}, selector string, inputValue map[string]interface{}) error {
	return p.exec(selector, inputValue, func(stmt *sqlx.NamedStmt, queryer *Queryer) error {
		if stmt == nil {
			return p.db.Get(dest, queryer.Sql)
		}
		return stmt.Get(dest, queryer.Value)
	})
}

func (p *GoBatis) QueryObjects(dest interface{}, selector string, inputValue map[string]interface{}) error {
	return p.exec(selector, inputValue, func(stmt *sqlx.NamedStmt, queryer *Queryer) error {
		if stmt == nil {
			return p.db.Select(dest, queryer.Sql)
		}
		return stmt.Select(dest, queryer.Value)
	})
}

func (p *GoBatis) update(selector string, inputValue map[string]interface{}) (int64, error) {
	var n int64 = 0
	err := p.exec(selector, inputValue, func(stmt *sqlx.NamedStmt, queryer *Queryer) error {
		var rs sql.Result
		var err error
		if stmt == nil {
			rs, err = p.db.Exec(queryer.Sql, queryer.Value)
		} else {
			rs, err = stmt.Exec(queryer.Value)
		}

		if err != nil {
			return err
		}

		if strings.ToUpper(queryer.Sql[0:6]) == "INSERT" {
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

func (p *GoBatis) getStmt(queryer *Queryer) (*sqlx.NamedStmt, error) {
	v, found := p.stmts.Load(queryer.Sql)
	if found {
		stmt := v.(*sqlx.NamedStmt)
		return stmt, nil
	}

	p.lock.Lock()
	v, found = p.stmts.Load(queryer.Sql)
	if found {
		stmt := v.(*sqlx.NamedStmt)
		p.lock.Unlock()
		return stmt, nil
	}

	var stmt *sqlx.NamedStmt
	var err error
	stmt, err = p.db.PrepareNamed(queryer.Sql)
	if err != nil {
		return nil, err
	}
	logger.Debug("prepared SQL => \n", stmt.QueryString)
	p.stmts.Store(queryer.Sql, stmt)
	p.lock.Unlock()

	return stmt, err
}
