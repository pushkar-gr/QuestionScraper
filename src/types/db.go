package types

import (
	"database/sql"
	"fmt"
	"os"
)

type DB struct {
	Username string `toml:"username"`
	DBName   string `toml:"dbname"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DB       *sql.DB
}

// insert platforms and topics to database
// input: Config
// output: error if any
func (db *DB) Init(config *Config) error {
	//get database password from environment variables
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return fmt.Errorf("DB_PASSWORD environment variable not set")
	}

	//construct database connection string
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%d sslmode=disable",
		config.Database.Username, config.Database.DBName, dbPassword, config.Database.Host, config.Database.Port)

	//connect to the database
	var err error
	db.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	//insert platforms
	for _, platform := range config.Platforms {
		//check if platform already exists in db
		var existingID int
		err := db.DB.QueryRow("SELECT id FROM platforms WHERE name = $1", platform.Name).Scan(&existingID)
		if err != nil {
			if err == sql.ErrNoRows {
				//platform does not exist
				_, err := db.DB.Exec("INSERT INTO platforms (name, website_url) VALUES ($1, $2)", platform.Name, platform.WebsiteURL)
				if err != nil {
					return fmt.Errorf("Error inserting platform %s: %v", platform.Name, err)
				} else {
					fmt.Printf("Inserted platform: %s\n", platform.Name)
				}
			} else {
				return fmt.Errorf("Error checking for existing platform %s: %v", platform.Name, err)
			}
		} else {
			fmt.Printf("Platform %s already exists with ID: %d\n", platform.Name, existingID)
		}
	}

	//insert topics
	for _, topic := range config.Topics {
		//check if topic already exists in db
		var existingID int
		err := db.DB.QueryRow("SELECT id FROM topics WHERE name = $1", topic.Name).Scan(&existingID)
		if err != nil {
			if err == sql.ErrNoRows {
				//topic does not exist
				_, err := db.DB.Exec("INSERT INTO topics (name, description) VALUES ($1, $2)", topic.Name, topic.Description)
				if err != nil {
					return fmt.Errorf("Error inserting topic %s: %v", topic.Name, err)
				} else {
					fmt.Printf("Inserted topic: %s\n", topic.Name)
				}
			} else {
				return fmt.Errorf("Error checking for existing topic %s: %v", topic.Name, err)
			}
		} else {
			fmt.Printf("Topic %s already exists with ID: %d\n", topic.Name, existingID)
		}
	}

	fmt.Println("Data population complete.\n\n")
	return nil
}
