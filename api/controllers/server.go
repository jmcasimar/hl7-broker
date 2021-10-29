package controllers

import (
	"fmt"
	"github.com/free-health/health24-gateway/conf"
	"github.com/free-health/health24-gateway/realtime"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func (s *Server) Init(db *gorm.DB) {
	s.DB = db
	s.Router = gin.New()
	s.initRoutes()
	go realtime.SocketIO.Serve()
}

func (s *Server) Run(config *conf.AppConfig) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Infof("running http server in %s", addr)
	log.Fatal(s.Router.Run(addr))
}
