package gobatis

type (
	selectorEntity struct {
		Name string
		Id   string
	}

	Queryer struct {
		Sql string
		// PREPARED (default), STATEMENT
		StatementType string
		Value         map[string]interface{}
	}

	DbConfig struct {
		Driver             string `json:"driver"`
		Dsn                string `json:"dsn"`
		MaxOpenConnections int    `json:"maxOpenConnections"`
		MaxIdleConnections int    `json:"maxIdleConnections"`
		MaxLifeTime        int    `json:"maxLifeTime"`
		MaxIdleTime        int    `json:"maxIdleTime"`
	}

	KeyValue struct {
		Key   string
		Value interface{}
	}
)
