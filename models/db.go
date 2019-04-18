package models

import (
	"fmt"
	"gitlab.com/paiduay/queq-hospital-api/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// InitDatabase - Init mysql database
func InitDatabase() error {
	connectionString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		config.Database.Username,
		config.Database.Password,
		config.Database.IP,
		config.Database.Port,
		config.Database.Database,
	)
	var err error
	db, err = gorm.Open("mssql", connectionString)
	return err
}
