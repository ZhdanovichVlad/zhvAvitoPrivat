package config

import (
	"fmt"
	"os"
)

type Config struct {
	//DatabasePort     string
	//DatabaseUser     string
	//DatabasePassword string
	//DatabaseName     string
	//DatabaseHost     string
	PgDSN      string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	var config *Config
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return config, fmt.Errorf("please set PG_DSG env")
	}

	pgDsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return config, fmt.Errorf("please set SERVER_PORT env")
	}

	config = &Config{
		//databasePort,
		//databaseUser,
		//databasePassword,
		//databaseName,
		//databaseHost,
		pgDsn,
		serverPort}

	return config, nil

}
