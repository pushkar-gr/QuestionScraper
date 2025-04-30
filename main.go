package main

import (
	"fmt"

	"github.com/pushkar-gr/QuestionScraper/src/geeksforgeeks"
	"github.com/pushkar-gr/QuestionScraper/src/hackerearth"
	"github.com/pushkar-gr/QuestionScraper/src/leetcode"
	"github.com/pushkar-gr/QuestionScraper/src/types"
)

func main() {
	config := new(types.Config)
	config.Update("config/config.toml")

	config.Database.Init(config)

	//call UpdateDB for all platforms
	for _, platform := range config.Platforms {
		fmt.Printf("Updating database for platform: %s\n", platform)

		switch platform.Name {
		case "GeeksForGeeks":
			geeksforgeeks.UpdateDB(config)
		case "LeetCode":
			leetcode.UpdateDB(config)
		case "HackerEarth":
			hackerearth.UpdateDB(config)
		default:
			fmt.Printf("Warning: No implementation found for platform %s\n", platform)
		}
	}
}
