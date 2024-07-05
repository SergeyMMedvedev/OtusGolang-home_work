package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type SchedulerConfig struct {
	Logger   LoggerConf
	Broker   BrokerConf
	Exchange ExchangeConf
	Storage  StorageConf
}

func NewSchedulerConfig() *SchedulerConfig {
	return &SchedulerConfig{}
}

func (s *SchedulerConfig) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *SchedulerConfig) Read(fpath string) (err error) {
	// read yaml file
	data, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, s)
	if err != nil {
		return err
	}
	return nil
}
