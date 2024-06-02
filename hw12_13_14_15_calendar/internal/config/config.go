package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type ServerConf struct {
	Host string
	Port int64
}

type LoggerConf struct {
	Level string
	// TODO
}

type StorageConf struct {
	Type string
	Psql PsqlConf
}

type PsqlConf struct {
	Host      string
	Port      int64
	User      string
	Password  string
	Dbname    string
	Sslmode   string
	Migration string
}

func NewConfig() Config {
	return Config{}
}

func (c *Config) Read(fpath string) (err error) {
	_, err = toml.DecodeFile(fpath, &c)
	return
}

func (c *Config) String() string {
	return fmt.Sprintf("%+v", *c)
}

// TODO
