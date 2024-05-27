package main

import "github.com/gin-gonic/gin"

func setUpRoutes(engine *gin.Engine) {
	engine.Use(gin.Recovery())
}
