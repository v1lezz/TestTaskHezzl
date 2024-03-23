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
	return fmt.Sprintf("")
}

type Config struct {
	dbCFG DBConfig
	port  int
}
