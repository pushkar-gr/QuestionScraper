package geeksforgeeks

import (
	"log"

	"github.com/pushkar-gr/QuestionScraper/src/types"
)

func UpdateDB(config *types.Config) {
	//iterate throught all topics
	for _, topic := range config.Topics {
		titles, err := GetTitleSlugs(topic.Name)
		if err != nil {
			log.Printf("Failed to get topic titles for %v from GeeksForGeeks: %v\n", topic.Name, err)
			continue
		}

		//iterate trhought all questions in GeeksForGeeks for topic
		for _, title := range titles {
			question, err := GetQuestion(title)
			if err != nil {
				log.Printf("Failed to get question %v for topic %v from GeeksForGeeks: %v", title, topic.Name, err)
				continue
			}

			//insert question to database
			err = config.Database.InsertQuestion(question)
			if err != nil {
				log.Printf("Failed to insert question %v for topic %v from GeeksForGeeks: %v", title, topic.Name, err)
			} else {
				log.Printf("Inserted question %v for topic %v from GeeksForGeeks", title, topic)
			}
		}
	}
}
