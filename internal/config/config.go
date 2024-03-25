package config

import "fmt"

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func (cfg DBConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

type Config struct {
	dbCFG DBConfig
	port  int
}
