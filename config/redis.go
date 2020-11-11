package config

type Redis struct {
	Addr     string `mapstructure:"redis" json:"addr" yaml:"addr"`
	Db       int    `mapstructure:"db" json:"db" yaml:"db"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}
