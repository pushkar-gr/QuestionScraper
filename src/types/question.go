package types

type DifficultyLevel string

const (
	Easy   DifficultyLevel = "Easy"
	Medium DifficultyLevel = "Medium"
	Hard   DifficultyLevel = "Hard"
)

type Question struct {
	Title       string
	Platform    string
	ExternalID  string
	Link        string
	Difficulty  DifficultyLevel
	Question    string
	Solution    string
	Explanation string
	Topics      []string
}
