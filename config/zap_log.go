package config

type ZapLog struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	Format        string `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Director      string `mapstructure:"director" json:"director"  yaml:"director"`
	LinkName      string `mapstructure:"link_name" json:"linkName" yaml:"link_name"`
	ShowLine      bool   `mapstructure:"show_line" json:"showLine" yaml:"showLine"`
	EncodeLevel   string `mapstructure:"encode_level" json:"encodeLevel" yaml:"encode_level"`
	StacktraceKey string `mapstructure:"stacktrace_key" json:"stacktraceKey" yaml:"stacktraceKey"`
	LogInConsole  bool   `mapstructure:"log_in_console" json:"logInConsole" yaml:"log_in_console"`
}
