package db

type DatabaseConfig struct {
	DBName string `json:"dbname"`
	DBPass string `json:"dbpass"`
	DBUser string `json:"dbuser"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}
