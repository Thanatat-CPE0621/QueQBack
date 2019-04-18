package config

// DatabaseConfiguration - Struct for database configuration
type DatabaseConfiguration struct {
	Production  databaseConfigBody `json:"production"`
	Development databaseConfigBody `json:"development"`
	Default     databaseConfigBody `json:"default"`
}

type databaseConfigBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
}

// Database - use to get current configuration
var Database databaseConfigBody

// SetDatabaseConfiguration - set database current configuration
func SetDatabaseConfiguration(readConfig DatabaseConfiguration, mode string) {
	Database = databaseConfigBody{}
	if mode == "prod" {
		Database = readConfig.Production
	} else if mode == "dev" {
		Database = readConfig.Development
	} else {
		Database = readConfig.Default
	}
}
