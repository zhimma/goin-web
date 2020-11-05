package config

type Jwt struct {
	JwtSecret Mysql `mapstructure:"jwt-secret" json:"jwtSecret" yaml:"jwt_secret"`
}
