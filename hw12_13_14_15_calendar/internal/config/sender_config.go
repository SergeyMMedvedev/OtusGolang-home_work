package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type SenderConfig struct {
	Logger  LoggerConf
	Broker  BrokerConf
	Queue   QueueConf
	Binding BindingConf
}

type QueueConf struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}
type BindingConf struct {
	QueueName string
	Exchange  string
	Key       string
	NoWait    bool
}

func NewSenderConfig() *SenderConfig {
	return &SenderConfig{}
}

func (s *SenderConfig) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *SenderConfig) Read(fpath string) (err error) {
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
