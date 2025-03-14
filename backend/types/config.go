package types

type (
	HttpConfig struct {
		ListenAddress string `json:"ListenAddress" yaml:"ListenAddress"`
		ListenPort    uint16 `json:"ListenPort" yaml:"ListenPort"`
		StaticRoot    string `json:"StaticRoot" yaml:"StaticRoot"`
	}
	DatabaseConfig struct {
		Rdb   *RdbConfig   `json:"Rdb" yaml:"Rdb"`
		Mongo *MongoConfig `json:"Mongo" yaml:"Mongo"`
	}
	RdbDriverType string
	RdbConfig     struct {
		DriverType            RdbDriverType `json:"DriverType" yaml:"DriverType"`
		Debug                 bool          `json:"Debug" yaml:"Debug"`
		DSN                   string        `json:"DSN" yaml:"DSN"`
		MaxIdleConnCount      int           `json:"MaxIdleConnCount" yaml:"MaxIdleConnCount"`
		MaxConnCount          int           `json:"MaxConnCount" yaml:"MaxConnCount"`
		ConnMaxIdleTimeSecond int           `json:"ConnMaxIdleTimeSecond" yaml:"ConnMaxIdleTimeSecond"`
		AutoMigrateLevel      string        `json:"AutoMigrateLevel" yaml:"AutoMigrateLevel"`
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
		Nas      NasConfig      `json:"Nas" yaml:"Nas"`
	}

	LogConfig struct {
		Level string `json:"Level" yaml:"Level"`
	}

	RealFilenamePolicy string
	NasConfig          struct {
		SimpleUploadRoot   string             `json:"SimpleUploadRoot" yaml:"SimpleUploadRoot"`
		RealFilenamePolicy RealFilenamePolicy `json:"RealFilenamePolicy" yaml:"RealFilenamePolicy"`
	}
)

const (
	DriverTypeSqlite   RdbDriverType = "sqlite"
	DriverTypePostgres RdbDriverType = "postgres"

	RFNP_Origin    RealFilenamePolicy = "origin"
	RFNP_Underline RealFilenamePolicy = "_"
	RFNP_UUID      RealFilenamePolicy = "uuid"
)

func (cfg *Config) SetDefault() {
	cfg.Http.SetDefault()
	cfg.Log.SetDefault()
	cfg.Database.Rdb.SetDefault()
	cfg.Database.Mongo.SetDefault()
	cfg.Nas.SetDefault()
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
	if cfg.StaticRoot == "" {
		cfg.StaticRoot = "./"
	}
}

func (c *RdbConfig) SetDefault() {
	if c == nil {
		return
	}
	if c.DriverType == "" {
		c.DriverType = DriverTypeSqlite
	}
	if c.DSN == "" {
		c.DSN = "file://localhost/rdb.db?mode=rwc"
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
	if c == nil {
		return
	}
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

func (c *NasConfig) SetDefault() {
	if c.RealFilenamePolicy == "" {
		c.RealFilenamePolicy = RFNP_UUID
	}
}
