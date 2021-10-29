package api

import (
	"github.com/free-health/health24-gateway/api/controllers"
	"github.com/free-health/health24-gateway/api/db"
	"github.com/free-health/health24-gateway/api/models"
	"github.com/free-health/health24-gateway/conf"
	"github.com/free-health/health24-gateway/realtime"
	log "github.com/sirupsen/logrus"
)

func Init(config conf.AppConfig, dbConfig conf.DBConfig) error {
	db.Init(dbConfig)
	models.Migrate()

	srv := controllers.Server{}
	realtime.StartSocketIO()
	srv.Init(db.DB)
	go srv.Run(&config)

	return nil
}

func Clean() {
	log.Infof("closing db")
	_ = db.DB.Close()
	//_ = realtime.SocketIO.Close()
}
