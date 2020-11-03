package config

type Base struct {
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	ZapLog ZapLog `mapstructure:"zap" json:"zap" yaml:"zap"`
	System System `mapstructure:"system" json:"zap" yaml:"zap"`
}
