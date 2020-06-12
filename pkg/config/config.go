package config

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	configFileName = "config.toml"
	setHandlerMode = "viper"
)

type Handler interface {
	ReadInConfig() error
	GetString(string) string
	GetInt(string) int
	GetDuration(string) time.Duration
	Set(string, interface{})
	SetConfigType(string)
	SetConfigFile(string)
	GetStringMapString(key string) map[string]string
	GetBool(string) bool
}

type Config struct {
	Path        string
	Name        string
	WatchConfig bool
	Handler
}

func (cfg *Config) setup(configName string) {
	if configName != "" {
		cfg.Name = configFileName
	}

	appPath, err := filepath.Abs(cfg.Name)
	//fmt.Println("appPath", appPath)
	if err != nil {
		panic(err)
	}
	cfg.Path = filepath.Dir(appPath)
	cfg.Path = strings.Replace(cfg.Path, "\\", "/", -1)

	cfg.handler()
}

func (cfg *Config) handler() {

	switch setHandlerMode {
	case "viper":
		cfg.Handler = viper.New()
	case "apolo":
		//cfg.Handler = apolo.New()
	default:
		panic("Fatal error setHandlerModle NOT FOUND \n")
	}

	cfg.Handler.Set("app_path", cfg.Path)
	cfg.Handler.SetConfigFile(filepath.Join(cfg.Path, configFileName))

	if err := cfg.Handler.ReadInConfig(); err != nil { // Handle errors reading the config file
		panic("Fatal error config file,err:" + err.Error() + " \n")
	}

}

var once sync.Once
var config Handler

func GetInstance() Handler {
	once.Do(func() {
		cfg := new(Config)
		cfg.setup(configFileName)
		config = cfg.Handler
	})
	return config
}
