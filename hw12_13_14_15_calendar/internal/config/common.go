package config

type LoggerConf struct {
	Level string
}

type BrokerConf struct {
	URI string
}

type ExchangeConf struct {
	Name     string
	Type     string
	Key      string
	Reliable bool
}
