package gormdialects

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func MySqlDB(uri string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(uri), &gorm.Config{})
}
