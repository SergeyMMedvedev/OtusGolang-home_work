package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type SenderConfig struct {
	Logger   LoggerConf
	Broker   BrokerConf
	Queue    QueueConf
	Binding  BindingConf
	Exchange ExchangeConf
	Consumer ConsumerConf
}

type ConsumerConf struct {
	Tag       string
	Lifetime  time.Duration
	NoAck     bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      map[string]interface{}
}

type QueueConf struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Key        string
	Args       map[string]interface{}
}
type BindingConf struct {
	QueueName string
	Exchange  string
	Key       string
	NoWait    bool
	Args      map[string]interface{}
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
