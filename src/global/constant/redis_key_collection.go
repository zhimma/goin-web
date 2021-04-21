package constant

const (

	/**
	 * @key_type    字符串
	 * @description 示例redis key
	 * @data        %s 值：用户id
	 */
	DEMO_REDIS_KEY_NAME string = "demo:redis_key:name:%s"

	// 管理后台登陆token key 用户id
	AdminManagerAccessToken string = "admin:manager:access_token:%v"

	// 管理后台登陆token key 用户id
	AdminManagerInfo string = "admin:manager:info:%v"

	/**
	 * @key_type    string
	 * @description 管理员用户信息
	 * @data        %s 用户client_id
	 */
	ClientAuthToken string = "passport:client:auth_token:%v"

	/**
	 * @key_type    string
	 * @description 管理员用户信息
	 * @data        %s 用户id
	 */
	ClientInfo = "passport:client:info:%v"

	/**
	 * @key_type    string
	 * @description 微信支付平台证书路径
	 * @data        key的路径
	 */
	WECHAT_OFFICIAL_CERT_PATH = "payment:wechat_official:cert_path"

	/**
	 * @key_type    string
	 * @description 微信支付平台证书编号
	 * @data        编号
	 */
	WECHAT_OFFICIAL_CERT_SERIAL_NUMBER = "payment:wechat_official:cert_serial_number"

	/**
	 * @key_type    string
	 * @description 微信支付平台证书编号
	 * @data        商户相关信息 %s=merchant_sn
	 */
	MERCHANT_INFO = "payment:merchant_info:%s"

	/**
	 * @key_type    string
	 * @description redis stream key 经销商信息变更同步
	 */
	LEGENDAGE_DEALER_SYNC = "legendage_dealer_sync"

	/**
	 * @key_type    string
	 * @description 客户端注册的回调url
	 * @data        商户相关信息 %s=client_id %s=module
	 */
	NOTIFY_CALLBACK_CLIENT = "notify:callback:client:%s:%s"
)
