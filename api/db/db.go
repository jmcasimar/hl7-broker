package db

import (
	"fmt"
	"github.com/free-health/health24-gateway/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

var DB *gorm.DB

func Init(config conf.DBConfig) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Name,
		config.Password,
	)
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}

	DB = db
}
