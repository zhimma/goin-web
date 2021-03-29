package config

type System struct {
	Env    string `mapstructure:"env" json:"env"`
	Addr   int    `mapstructure:"addr" json:"addr"`
	DbType string `mapstructure:"db_type" json:"dbType"`
}
