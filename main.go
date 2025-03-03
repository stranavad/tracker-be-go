package main

import (
	"tracker/db"
	"tracker/session"
	"tracker/tracker"
	"tracker/types"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length", "Content-Type", "Authorization"},
	}))


	serviceConfig := types.ServiceConfig {
		DB: db.GetDb(),
	}

	tracker.RegisterRoutes(r, serviceConfig)
	session.RegisterRoutes(r, serviceConfig)
	r.Run() // listen and serve on 0.0.0.0:8080
}
