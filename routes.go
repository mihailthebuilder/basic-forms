package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setUpRoutes(r *gin.Engine) {
	r.GET("/user/:id/submissions", getSubmissions)
}

func getSubmissions(c *gin.Context) {
	c.String(http.StatusOK, "Return")
}
