package config

import (
	"time"
)

type Config struct {
	// common
	App   AppConfig
	Log   LogConfig
	MySQL struct {
		Default  MySQLConfig
		Material MySQLConfig
	}

	Redis struct {
		Default RedisConfig
	}
	GCP               GCPConfig
	YoutuberCrawler   YoutuberCrawler
	EmailAPI          EmailAPI
	OAuth             OAuth
	ContentModeration ContentModeration
	Docusign          Docusign //添加docusign
	WhiteList         []string
}

// AppConfig
type AppConfig struct {
	Name       string
	InstantKey string
	RunMode    string
	UseSwagger bool
	Addr       string
	Url        string
	JwtSecret  string
	JwtExpire  int
}

// LogConfig
type LogConfig struct {
	Writers          string
	LoggerLevel      string
	LoggerFile       string
	LogFormatText    bool
	LogRollingPolicy string
	LogRotateDate    int
	LogRotateSize    int
	LogBackupCount   int
}

// MySQLConfig
type MySQLConfig struct {
	Name            string
	Addr            string
	UserName        string
	Password        string
	ShowLog         bool
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime int
}

// RedisConfig
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	DialTimeOut  time.Duration
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	PoolSize     int
	PoolTimeOut  time.Duration
}

// GCPConfig
type GCPConfig struct {
	CDN     string
	KeyName string
	Key     string
}
type EmailAPI struct {
	Domain    string
	Uri       string
	GameId    string
	Sigkey    string
	ChannelId int
	Source    int
}

type OAuth struct {
	Youtube Youtube
}

type Youtube struct {
	AuthUrl          string
	TokenUrl         string
	RedirectUri      string
	ClientId         string
	ClientSecret     string
	LoginRedirectUri string
	YoutubeMinFans   int
}

type YoutuberCrawler struct {
	Default struct {
		Domain string
	}
}

type ContentModeration struct {
	Domain              string
	Uri                 string
	Appkey              string
	Appid               string
	Gameid              string
	BusiheadAppkey      string
	BusiheadServicename string
}

//docusign配置项
type Docusign struct {
	Env          string //环境
	PublicKey    string //公钥
	PrivateKey   string //私钥
	ClientId     string //
	RedirectUri  string
	WebhookUri   string //webHook地址
	Sub          string
	TemplateId   map[string]string
	ClientRegion string
	ConnectKey   string
}
