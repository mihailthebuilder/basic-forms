package main

import (
	"basic-forms/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type User struct {
	InternalId string `json:"internalId"`
	ExternalId string `json:"externalId"`
}

func setUpRoutes(r *gin.Engine) {
	r.POST("/users", createUser)
	r.POST("/submit/:externalUserId", postSubmission)
	r.GET("/users/:internalUserId/origins/:origin", getSubmissionsForOrigin)
}

func createUser(c *gin.Context) {
	internalUserId := uuid.New()

	externalId, err := encryption.Encrypt(internalUserId.String())
	if err != nil {
		logger.Error("can't encrypt internal ID: ", internalUserId, " error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = datastore.NewUser(internalUserId.String())
	if err != nil {
		logger.Error("can't create new user: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type User struct {
		InternalId string `json:"internalId"`
		ExternalId string `json:"externalId"`
	}

	c.JSON(http.StatusCreated, User{InternalId: internalUserId.String(), ExternalId: externalId})
}

func postSubmission(c *gin.Context) {
	externalUserId := c.Param("externalUserId")
	origin := c.Request.Header.Get("Origin")

	body, err := c.GetRawData()
	if err != nil {
		logger.Error("can't fetch body: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	internalUserId, err := encryption.Decrypt(externalUserId)
	if err != nil {
		logger.Error("can't decrypt externalUserId: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = datastore.AddSubmission(internalUserId, origin, body)
	if err != nil {
		logger.Error("can't add submission: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}

func getSubmissionsForOrigin(c *gin.Context) {
	internalUserId := c.Param("internalUserId")
	origin := c.Param("origin")

	content, err := datastore.GetSubmissions(internalUserId, origin)
	if err != nil {
		logger.Error("can't fetch file contents: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "text/plain", content)
}
