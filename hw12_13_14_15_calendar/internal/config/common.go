package config

type LoggerConf struct {
	Level string
}

type BrokerConf struct {
	URI string
}

type ExchangeConf struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Key        string
	Reliable   bool
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
