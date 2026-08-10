package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"
)

type Config struct {
	Env           string         `yaml:"env"`
	Server        *Server        `yaml:"server"`
	DB            *DB            `yaml:"db"`
	Redis         *Redis         `yaml:"redis"`
	JWT           *JWTConfig     `yaml:"jwt"`
	Zap           *ZapConfig     `yaml:"zap"`
	Swagger       *SwaggerConfig `yaml:"swagger"`
	AuthSkipPaths []string       `yaml:"authSkipPaths"`
}

type SwaggerConfig struct {
	Enabled bool     `yaml:"enabled"`
	Host    string   `yaml:"host"`
	Scheme  []string `yaml:"scheme"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type DB struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Name            string `yaml:"name"`
	Charset         string `yaml:"charset"`
	Collation       string `yaml:"collation"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime string `yaml:"connMaxLifetime"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	SecretKey string `yaml:"secretKey"`
	Expire    int    `yaml:"expire"`
}

type ZapConfig struct {
	OutputPaths      string `yaml:"outputPaths" json:"outputPaths"`
	ErrorOutputPaths string `yaml:"errorOutputPaths" json:"errorOutputPaths"`
}

var Configs *Config

func ReadConfig2() {
	file, err := os.Open("./configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	contents, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	Configs = &Config{}
	if err := yaml.Unmarshal(contents, Configs); err != nil {
		log.Fatal(err)
	}
}

func ReadConfigByViper() {
	viper := viper.New()
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")

	//  系统目录（生产环境）
	// viper.AddConfigPath("/etc/myapp/")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(fmt.Errorf("读取配置失败: %w", err))
	}
	Configs = &Config{}
	if err := viper.Unmarshal(Configs); err != nil {
		fmt.Println(fmt.Errorf("解析配置失败: %w", err))
	}

}

// 参考
// https://www.cnblogs.com/jiujuan/p/13799976.html
// viper 读取配置文件的优先顺序，从高到低，如下： Viper 配置key是不区分大小写的

// 1显式设置的Set函数
//viper.Set("mysql.info", "this is mysql info")

// 2命令行参数

// 3环境变量

// 4配置文件
// viper.SetConfigName("config") // 配置文件的文件名，没有扩展名，如 .yaml, .toml 这样的扩展名
// viper.SetConfigType("yaml")  // 设置扩展名。在这里设置文件的扩展名。另外，如果配置文件的名称没有扩展名，则需要配置这个选项
// viper.AddConfigPath("/etc/appname/") // 查找配置文件所在路径
// viper.AddConfigPath("$HOME/.appname") // 多次调用AddConfigPath，可以添加多个搜索路径
// viper.AddConfigPath(".")             // 还可以在工作目录中搜索配置文件
// err := viper.ReadInConfig()       // 搜索并读取配置文件
// if err != nil { // 处理错误
//   panic(fmt.Errorf("Fatal error config file: %s \n", err))
// }

// 5远程k-v 存储系统，如consul，etcd等

// 6默认值
// viper.SetDefault("ContentDir", "content")
