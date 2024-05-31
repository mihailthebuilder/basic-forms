package main

import (
	"basic-forms/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setUpRoutes(r *gin.Engine) {
	r.POST("/users", createUser)
	r.POST("/users/:userId/submit", postSubmission)
	r.GET("/users/:userId/origins/:origin", getSubmissionsForOrigin)
}

func createUser(c *gin.Context) {
	user, err := datastore.NewUser()

	if err != nil {
		logger.Error("can't create new user: ", err)
	}

	c.JSON(http.StatusCreated, user)
}

func postSubmission(c *gin.Context) {
	userId := c.Param("userId")
	origin := c.Request.Header.Get("Origin")

	body, err := c.GetRawData()
	if err != nil {
		logger.Error("can't fetch body: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = datastore.AddSubmission(userId, origin, body)
	if err != nil {
		logger.Error("can't add submission: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func getSubmissionsForOrigin(c *gin.Context) {
	userId := c.Param("userId")
	origin := c.Param("origin")

	content, err := datastore.GetSubmissions(userId, origin)
	if err != nil {
		logger.Error("can't fetch file contents: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "text/plain", content)
}
