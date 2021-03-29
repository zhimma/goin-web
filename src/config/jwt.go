package config

type Jwt struct {
	JwtSecret string `mapstructure:"jwt_secret" json:"jwtSecret" yaml:"jwt_secret"`
	JwtTtl    int64  `mapstructure:"jwt_ttl" json:"jwtTtl" yaml:"jwt_ttl"`
}
