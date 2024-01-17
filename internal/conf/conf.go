package conf

import (
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Db     Database   `yaml:"mysql" json:"db"`
	Server HttpServer `yaml:"server" json:"server"`
}

type Database struct {
	UserName     string        `yaml:"user_name" json:"user_name"`
	Password     string        `yaml:"password" json:"password"`
	Db           string        `yaml:"db" json:"db"`
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	Timeout      time.Duration `yaml:"timeout" json:"timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
}

type HttpServer struct {
	Listen string `yaml:"listen" json:"listen"`
	Port   int    `yaml:"port" json:"port"`
}

func MustLoad(config string) *Config {
	f, err := os.Open(config)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	c := &Config{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
	return c
}
