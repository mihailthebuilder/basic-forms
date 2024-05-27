package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setUpRoutes(r *gin.Engine) {
	r.GET("/user/:id/submissions", getSubmissions)
	r.POST("/user/:id/submit", postSubmission)
}

func getSubmissions(c *gin.Context) {
	c.String(http.StatusOK, "Return")
}

func postSubmission(c *gin.Context) {
	c.String(http.StatusOK, "Return")
}
