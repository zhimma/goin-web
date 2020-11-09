package config

type Jwt struct {
	JwtSecret        string `mapstructure:"jwt_secret" json:"jwtSecret" yaml:"jwt_secret"`
	JwtTtl           int64  `mapstructure:"jwt_ttl" json:"jwtTtl" yaml:"jwt_ttl"`
	JwtRefreshSecret string `mapstructure:"jwt_refresh_secret" json:"jwtRefreshSecret" yaml:"jwt_refresh_secret"`
	JwtRefreshTtl    int64  `mapstructure:"jwt_refresh_ttl" json:"JwtRefreshTtl" yaml:"jwt_refresh_ttl"`
}
