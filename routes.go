package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setUpRoutes(r *gin.Engine) {
	r.GET("/users/:userId/submissions", getSubmissions)
	r.POST("/users/:userId/submit", postSubmission)
}

func getSubmissions(c *gin.Context) {
	c.String(http.StatusOK, "Return")
}

func postSubmission(c *gin.Context) {
	userId := c.Param("userId")
	host := c.Request.Header.Get("Origin")

	log.Println("request", userId, host)

	c.String(http.StatusOK, "Return")
}
