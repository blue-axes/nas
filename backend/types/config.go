package types

type (
	HttpConfig struct {
		ListenAddress string `json:"ListenAddress" yaml:"ListenAddress"`
		ListenPort    uint16 `json:"ListenPort" yaml:"ListenPort"`
	}
	DatabaseConfig struct {
		Postgres *PostgresConfig `json:"Postgres" yaml:"Postgres"`
		Mongo    *MongoConfig    `json:"Mongo" yaml:"Mongo"`
	}
	PostgresConfig struct {
		Debug                 bool   `json:"Debug" yaml:"Debug"`
		Address               string `json:"Address" yaml:"Address"`
		Port                  uint16 `json:"Port" yaml:"Port"`
		Username              string `json:"Username" yaml:"Username"`
		Password              string `json:"Password" yaml:"Password"`
		Database              string `json:"Database" yaml:"Database"`
		MaxIdleConnCount      int    `json:"MaxIdleConnCount" yaml:"MaxIdleConnCount"`
		MaxConnCount          int    `json:"MaxConnCount" yaml:"MaxConnCount"`
		ConnMaxIdleTimeSecond int    `json:"ConnMaxIdleTimeSecond" yaml:"ConnMaxIdleTimeSecond"`
		AutoMigrateLevel      string `json:"AutoMigrateLevel" yaml:"AutoMigrateLevel"`
	}

	MongoConfig struct {
		Debug                 bool   `json:"Debug" yaml:"Debug"`
		Address               string `json:"Address" yaml:"Address"`
		Port                  uint16 `json:"Port" yaml:"Port"`
		Username              string `json:"Username" yaml:"Username"`
		Password              string `json:"Password" yaml:"Password"`
		Database              string `json:"Database" yaml:"Database"`
		MaxIdleConnCount      int    `json:"MaxIdleConnCount" yaml:"MaxIdleConnCount"`
		MaxConnCount          int    `json:"MaxConnCount" yaml:"MaxConnCount"`
		ConnMaxIdleTimeSecond int    `json:"ConnMaxIdleTimeSecond" yaml:"ConnMaxIdleTimeSecond"`
		AutoMigrateLevel      string `json:"AutoMigrateLevel" yaml:"AutoMigrateLevel"`
	}

	Config struct {
		Http     HttpConfig     `json:"Http" yaml:"Http"`
		Log      LogConfig      `json:"Log" yaml:"Log"`
		Database DatabaseConfig `json:"Database" yaml:"Database"`
	}

	LogConfig struct {
		Level string `json:"Level" yaml:"Level"`
	}
)

func (cfg *Config) SetDefault() {
	cfg.Http.SetDefault()
	cfg.Log.SetDefault()
	cfg.Database.Postgres.SetDefault()
}

func (cfg *LogConfig) SetDefault() {
	if cfg.Level == "" {
		cfg.Level = "info"
	}
}

func (cfg *HttpConfig) SetDefault() {
	if cfg.ListenAddress == "" {
		cfg.ListenAddress = "0.0.0.0"
	}
	if cfg.ListenPort <= 0 {
		cfg.ListenPort = 80
	}
}

func (c *PostgresConfig) SetDefault() {
	if c.Address == "" {
		c.Address = "127.0.0.1"
	}
	if c.Port <= 0 {
		c.Port = 5432
	}
	if c.Username == "" {
		c.Username = "postgres"
	}
	if c.MaxConnCount <= 0 {
		c.MaxConnCount = 5
	}
	if c.MaxIdleConnCount <= 0 {
		c.MaxIdleConnCount = 2
	}
	if c.ConnMaxIdleTimeSecond <= 0 {
		c.ConnMaxIdleTimeSecond = 60 * 5
	}
}

func (c *MongoConfig) SetDefault() {
	if c.Address == "" {
		c.Address = "127.0.0.1"
	}
	if c.Port <= 0 {
		c.Port = 27017
	}
	if c.Username == "" {
		c.Username = "admin"
	}
	if c.MaxConnCount <= 0 {
		c.MaxConnCount = 5
	}
	if c.MaxIdleConnCount <= 0 {
		c.MaxIdleConnCount = 2
	}
	if c.ConnMaxIdleTimeSecond <= 0 {
		c.ConnMaxIdleTimeSecond = 60 * 5
	}
}
