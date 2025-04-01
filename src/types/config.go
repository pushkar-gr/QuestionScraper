package types

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database  DB `toml:"database"`
	Platforms []struct {
		Name       string `toml:"name"`
		WebsiteURL string `toml:"website_url"`
	} `toml:"platforms"`
	Topics []struct {
		Name        string `toml:"name"`
		Description string `toml:"description"`
	}
}

// update Config using data from config file and environment variables
// input: config file path
// output: error if any
func (config *Config) Update(path string) error {
	//update Config using config file
	if err := config.UpdatePath(path); err != nil {
		return err
	}
	//overwrite Config if data found in environment variable
	return config.UpdateENV()
}

// update Config using data from config file
// input: config file path
// output: error if any
func (config *Config) UpdatePath(path string) error {
	//read config file and upadte Config
	_, err := toml.DecodeFile(path, config)
	return err
}

// update config using data from config file
// input: config file path
// output: error if any
func (config *Config) UpdateENV() error {
	//get username from environment variables
	dbUsername := os.Getenv("DB_USERNAME")
	if dbUsername != "" {
		config.Database.Username = dbUsername
	}

	//get database name from environment variables
	dbName := os.Getenv("DB_NAME")
	if dbName != "" {
		config.Database.DBName = dbName
	}

	//get host from environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		config.Database.Host = dbHost
	}

	//get port from environment variables
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr != "" {
		dbPort, err := strconv.Atoi(dbPortStr)
		if err != nil {
			return fmt.Errorf("Invalid DB_PORT environment variable: %v", err)
		}
		config.Database.Port = dbPort
	}

	return nil
}
