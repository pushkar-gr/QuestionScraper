package types

type DB struct {
	Username string `toml:"username"`
	DBName   string `toml:"dbname"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
}
