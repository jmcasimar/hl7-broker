package models

import (
	"github.com/free-health/health24-gateway/api/db"
	log "github.com/sirupsen/logrus"
)

func Migrate() {
	log.Infof("migrating models")
	db.DB.AutoMigrate(
		&User{},
		&Patient{},
		&PhysiologicalAlarm{},
		&Monitor{},
	)
}
