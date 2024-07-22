package config

var Config configStruct

type configStruct struct {
	Env struct {
		Profiles string `yaml:"profiles"`
	} `yaml:"env"`
	Mysql struct {
		Address       []string `yaml:"address"`
		Username      string   `yaml:"username"`
		Password      string   `yaml:"password"`
		Database      string   `yaml:"database"`
		MaxOpenConn   int      `yaml:"maxOpenConn"`
		MaxIdleConn   int      `yaml:"maxIdleConn"`
		MaxLifeTime   int      `yaml:"maxLifeTime"`
		LogLevel      int      `yaml:"logLevel"`
		SlowThreshold int      `yaml:"slowThreshold"`
	} `yaml:"mysql"`

	Mongo struct {
		Uri         string   `yaml:"uri"`
		Address     []string `yaml:"address"`
		Database    string   `yaml:"database"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
		MaxPoolSize int      `yaml:"maxPoolSize"`
	} `yaml:"mongo"`

	Redis struct {
		ClusterMode bool     `yaml:"clusterMode"`
		Address     []string `yaml:"address"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
		DB          int      `yaml:"db"`
		PoolSize    int      `yaml:"poolSize"`
	} `yaml:"redis"`

	FileStore struct {
		DriverType      string `yaml:"driverType"`
		Endpoint        string `yaml:"endpoint"`
		AccessKeyID     string `yaml:"accessKeyId"`
		AccessKeySecret string `yaml:"accessKeySecret"`
		BucketName      string `yaml:"bucketName"`
		Domain          string `yaml:"domain"`
	} `yaml:"fileStore"`

	Sms struct {
		DriverType      string `yaml:"driverType"`
		Endpoint        string `yaml:"endpoint"`
		AccessKeyID     string `yaml:"accessKeyId"`
		AccessKeySecret string `yaml:"accessKeySecret"`
	} `yaml:"sms"`

	Mail struct {
		DriverType string `yaml:"driverType"`
		From       string `yaml:"from"`
		Account    string `yaml:"account"`
		Password   string `yaml:"password"`
		SmtpHost   string `yaml:"smtpHost"`
		SmtpPort   int    `yaml:"smtpPort"`
	} `yaml:"mail"`

	WxOpenPlatform struct {
		AppID         string `yaml:"appId"`
		AppSecret     string `yaml:"appSecret"`
		MessageToken  string `yaml:"messageToken"`
		MessageAesKey string `yaml:"messageAesKey"`
	} `yaml:"wxOpenPlatform"`

	WxOfficialAccount struct {
		AppID string `yaml:"appId"`
	} `yaml:"wxOfficialAccount"`

	Log struct {
		StorageLocation     string `yaml:"storageLocation"`
		RotationTime        uint   `yaml:"rotationTime"`
		RemainRotationCount uint   `yaml:"remainRotationCount"`
		RemainLogLevel      int    `yaml:"remainLogLevel"`
		IsStdout            bool   `yaml:"isStdout"`
		IsJson              bool   `yaml:"isJson"`
		WithStack           bool   `yaml:"withStack"`
	} `yaml:"log"`

	Prometheus struct {
		Enable              bool   `yaml:"enable"`
		PrometheusUrl       string `yaml:"prometheusUrl"`
		AdminPrometheusPort []int  `yaml:"adminPrometheusPort"`
		AppPrometheusPort   []int  `yaml:"appPrometheusPort"`
	} `yaml:"prometheus"`

	Secret      string `yaml:"secret"`
	TokenPolicy struct {
		Expire int64 `yaml:"expire"`
	} `yaml:"tokenPolicy"`

	SuperAdmin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		RealName string `yaml:"realName"`
		Phone    string `yaml:"phone"`
	} `yaml:"superAdmin"`

	AdminApi struct {
		Port     []int  `yaml:"port"`
		ListenIP string `yaml:"listenIP"`
	} `yaml:"adminApi"`

	AppApi struct {
		Port     []int  `yaml:"port"`
		ListenIP string `yaml:"listenIP"`
	} `yaml:"appApi"`
}
