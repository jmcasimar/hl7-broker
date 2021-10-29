package main

import (
	"github.com/free-health/health24-gateway/api"
	"github.com/free-health/health24-gateway/conf"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"os/signal"
	"syscall"

	"github.com/free-health/health24-gateway/discovery"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	gin.SetMode("release")

	err := godotenv.Load()
	if err != nil {
		log.Error("error loading .env file")
	}

	var config conf.AppConfig
	err = envconfig.Process("app", &config)
	if err != nil {
		log.Fatal(err)
	}

	var dbConf conf.DBConfig
	err = envconfig.Process("db", &dbConf)
	if err != nil {
		log.Fatal(err)
	}

	err = api.Init(config, dbConf)
	if err != nil {
		log.Fatal(err)
	}

	discovery.Start(config.UdpDiscoveryHost)

	// wait for Ctrl+C
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)
	<-quit
	log.Infof("exiting program on SIGTERM")
	api.Clean()
	discovery.Stop()
	os.Exit(1)
}
