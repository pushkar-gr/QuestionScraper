package types

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
)

type DB struct {
	Username string `toml:"username"`
	DBName   string `toml:"dbname"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	DB       *sql.DB
}

// insert platforms and topics to database
// input config
// output error if any
func (db *DB) Populate(config *Config) error {
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

	return nil
}

// insert question into db
// input: Question
// output: error if any
func (db *DB) InsertQuestion(question Question) error {
	//begin transaction
	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("transaction begin error: %v", err)
	}
	defer tx.Rollback()

	//get or create platform
	var platformID int
	err = tx.QueryRow(`
			INSERT INTO platforms (name) 
			VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id`,
		question.Platform,
	).Scan(&platformID)

	if err != nil {
		return fmt.Errorf("platform error: %v", err)
	}

	// insert question
	var questionID int
	err = tx.QueryRow(`
     INSERT INTO questions
     (title, platform_id, external_id, link, difficulty, solution, explanation)
     VALUES ($1, $2, $3, $4, $5, $6, $7)
     RETURNING id`,
		question.Title,
		platformID,
		question.ExternalID,
		question.Link,
		question.Difficulty,
		question.Solution,
		question.Explanation,
	).Scan(&questionID)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // Unique constraint violation
				switch pqErr.Constraint {
				case "uq_platform_external":
					return fmt.Errorf("duplicate external ID for this platform")
				case "questions_link_key":
					return fmt.Errorf("duplicate question link")
				case "uq_platform_title":
					return fmt.Errorf("duplicate title for this platform")
				}
			}
		}
		return fmt.Errorf("question insert error: %v", err)
	}

	//process topics
	seenTopics := make(map[string]struct{})
	for _, topic := range question.Topics {
		seenTopics[topic] = struct{}{}
	}

	for topicName := range seenTopics {
		var topicID int
		err = tx.QueryRow(`
				INSERT INTO topics (name)
				VALUES ($1)
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
				RETURNING id`,
			topicName,
		).Scan(&topicID)

		if err != nil {
			return fmt.Errorf("topic error: %v", err)
		}

		_, err = tx.Exec(`
				INSERT INTO question_topic (question_id, topic_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING`,
			questionID,
			topicID,
		)

		if err != nil {
			return fmt.Errorf("topic mapping error: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %v", err)
	}

	return nil
}

// closes db
func (db *DB) CloseDB() {
	db.DB.Close()
}
