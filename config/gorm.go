package config

type Mysql struct {
	Host               string `mapstructure:"host" json:"host" yaml:"host"`
	Config             string `mapstructure:"config" json:"config" yaml:"config"`
	Dbname             string `mapstructure:"db_name" json:"dbname" yaml:"db_name"`
	Username           string `mapstructure:"username" json:"username" yaml:"username"`
	Password           string `mapstructure:"password" json:"password" yaml:"password"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections" json:"maxIdleConnections" yaml:"max_idle_connections"`
	MaxOpenConnections int    `mapstructure:"max_open_connections" json:"maxOpenConnections" yaml:"max_idle_connections"`
	LogMode            bool   `mapstructure:"log_mode" json:"logMode" yaml:"logger_mode"`
}
