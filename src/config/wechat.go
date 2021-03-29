package config

type Wechat struct {
	AppId         string `mapstructure:"app_id" json:"app_id" yaml:"app_id"`
	MiniAppId     string `mapstructure:"mini_app_id" json:"mini_app_id" yaml:"mini_app_id"`
	MchId         string `mapstructure:"mch_id" json:"mch_id" yaml:"mch_id"`
	NotifyUrl     string `mapstructure:"notify_url" json:"notify_url" yaml:"notify_url"`
	RefundUrl     string `mapstructure:"refund_url" json:"refund_url" yaml:"refund_url"`
	V2key         string `mapstructure:"v2_key" json:"v2_key" yaml:"v2_key"`
	V3Key         string `mapstructure:"v3_key" json:"v3_key" yaml:"v3_key"`
	Mode          string `mapstructure:"mode" json:"mode" yaml:"mode"`
	SerialNumber  string `mapstructure:"serial_number" json:"serial_number" yaml:"serial_number"`
	WechatCertPub string `mapstructure:"wechat_cert_pub" json:"wechat_cert_pub" yaml:"wechat_cert_pub"`
	SubMiniAppId  string `mapstructure:"sub_mini_app_id" json:"sub_mini_app_id" yaml:"sub_mini_app_id"`
}
