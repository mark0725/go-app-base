package db

type SqlBuildResult struct {
	Sql    string
	Params []interface{}
}

type DatabaseConfig struct {
	Type    string `json:"type,omitempty"`
	Dsn     string `json:"dsn,omitempty"`
	DBName  string `json:"dbname,omitempty"`
	DBPass  string `json:"dbpass,omitempty"`
	DBUser  string `json:"dbuser,omitempty"`
	Host    string `json:"host,omitempty"`
	Port    int    `json:"port,omitempty"`
	Driver  string `json:"driver,omitempty"`
	Options string `json:"options,omitempty"`
}

type PageQueryResult struct {
	TotalElements    int              `json:"totalElements"`
	TotalPages       int              `json:"totalPages"`
	Last             bool             `json:"last"`
	NumberOfElements int              `json:"numberOfElements"`
	First            bool             `json:"first"`
	Size             int              `json:"size"`
	Number           int              `json:"number"`
	Content          []map[string]any `json:"content"`
}
