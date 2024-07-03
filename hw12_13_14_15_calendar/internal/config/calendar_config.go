package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type CalendarConfig struct {
	Logger      LoggerConf
	Storage     StorageConf
	GRPCGateWay GRPCGateWayConf
	GRPCServer  GRPCServerConf
}

type GRPCGateWayConf struct {
	Host string
	Port int64
}

type GRPCServerConf struct {
	Host string
	Port int64
}

type StorageConf struct {
	Type string
	Psql PsqlConf
}

type PsqlConf struct {
	Host          string
	Port          int64
	User          string
	Password      string
	Dbname        string
	Sslmode       string
	MigrationDir  string `yaml:"migration_dir"`
	ExecMigration bool   `yaml:"exec_migration"`
}

func NewCalendarConfig() CalendarConfig {
	return CalendarConfig{}
}

func (c *CalendarConfig) Read(fpath string) (err error) {
	// read yaml file
	data, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *CalendarConfig) String() string {
	return fmt.Sprintf("%+v", *c)
}

// TODO
