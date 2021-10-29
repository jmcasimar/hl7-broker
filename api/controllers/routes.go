package controllers

import (
	"github.com/free-health/health24-gateway/api/middlewares"
	"github.com/free-health/health24-gateway/realtime"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func (s *Server) initRoutes() {
	s.Router.Use(middlewares.GinMiddlewareCORS("http://localhost:8081"))

	api := s.Router.Group("/api")
	{
		api.GET("/monitors", s.GetAllMonitors)
	}

	s.Router.GET("/socket.io/*any", gin.WrapH(realtime.SocketIO))
	s.Router.POST("/socket.io/*any", gin.WrapH(realtime.SocketIO))
	s.Router.Use(static.Serve("/", static.LocalFile("./public", true)))
}
