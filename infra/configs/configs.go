package configs

import (
	"flag"
	"forge/constant"
	"forge/pkg/log/zlog"
	"github.com/spf13/viper"
	"time"
)

type IConfig interface {
	GetRedisConfig() RedisConfig
	GetDBConfig() DBConfig
	GetAppConfig() ApplicationConfig
	GetLoggerConfig() LoggerConfig
}

var (
	conf = new(config)
)

func Config() IConfig {
	return conf
}
func MustInit(path string) {
	mustInit(path)
	return
}

func (c *config) GetRedisConfig() RedisConfig {
	return c.RedisConfig

}

func (c *config) GetDBConfig() DBConfig {
	return c.DBConfig
}

func (c *config) GetAppConfig() ApplicationConfig {
	return c.AppConfig
}

func (c *config) GetLoggerConfig() LoggerConfig {
	return c.LogConfig
}

func mustInit(path string) *config {
	// 初始化时间为东八区的时间
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	// 默认配置文件路径
	var configPath string
	flag.StringVar(&configPath, "c", path+constant.DEFAULT_CONFIG_FILE_PATH, "配置文件绝对路径或相对路径")
	flag.Parse()
	zlog.Infof("配置文件路径为 %s", configPath)
	// 初始化配置文件
	viper.SetConfigFile(configPath)
	viper.WatchConfig()
	// 观察配置文件变动
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	zlog.Warnf("配置文件发生变化")
	//	if err := viper.Unmarshal(&configs.Conf); err != nil {
	//		zlog.Errorf("无法反序列化配置文件 %v", err)
	//	}
	//	zlog.Debugf("%+v", configs.Conf)
	//
	//	Eve()
	//	Init()
	//})
	// 将配置文件读入 viper
	if err := viper.ReadInConfig(); err != nil {
		zlog.Panicf("无法读取配置文件 err: %v", err)
	}
	_config := config{}
	// 解析到变量中
	if err := viper.Unmarshal(&_config); err != nil {
		zlog.Panicf("无法解析配置文件 err: %v", err)
	}
	zlog.Debugf("配置文件为 ： %+v", _config)
	conf = &_config
	return conf

}

type config struct {
	AppConfig   ApplicationConfig `mapstructure:"app"`
	LogConfig   LoggerConfig      `mapstructure:"log"`
	DBConfig    DBConfig          `mapstructure:"database"`
	RedisConfig RedisConfig       `mapstructure:"redis"`
}

type ApplicationConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Env         string `mapstructure:"env"`
	LogfilePath string `mapstructure:"logfilePath"`
}
type LoggerConfig struct {
	Level    int8   `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Director string `mapstructure:"director"`
	ShowLine bool   `mapstructure:"show-line"`
}

type DBConfig struct {
	Driver      string `mapstructure:"driver"`
	AutoMigrate bool   `mapstructure:"migrate"`
	Dsn         string `mapstructure:"dsn"`
}
type RedisConfig struct {
	Enable   bool   `mapstructure:"enable"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type KafkaConfig struct {
	host string `mapstructure:"host"`
	port int    `mapstructure:"port"`
}
