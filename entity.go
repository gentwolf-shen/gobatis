package gobatis

type (
	selectorEntity struct {
		Name string
		Id   string
	}

	DbConfig struct {
		Driver             string `json:"driver"`
		Dsn                string `json:"dsn"`
		MaxOpenConnections int    `json:"maxOpenConnections"`
		MaxIdleConnections int    `json:"maxIdleConnections"`
		MaxLifeTime        int    `json:"maxLifeTime"`
		MaxIdleTime        int    `json:"maxIdleTime"`
	}
)
