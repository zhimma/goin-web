package config

type Casbin struct {
	ModelPath string `mapstructure:"model_path" json:"model_path" yaml:"model_path"`
}
