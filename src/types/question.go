package types

type Question struct {
	Title       string
	Platform    string
	ExternalID  string
	Link        string
	Difficulty  string
	Question    string
	Solution    string
	Explanation string
	Topics      []string
}
