package db

type DatabaseConfig struct {
	Type    string `json:"type"`
	DBName  string `json:"dbname"`
	DBPass  string `json:"dbpass"`
	DBUser  string `json:"dbuser"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Driver  string `json:"driver"`
	Options string `json:"options"`
}
