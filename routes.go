package main

import (
	"basic-forms/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setUpRoutes(r *gin.Engine) {
	r.POST("/users", createUser)
	r.POST("/users/:userId/submit", postSubmission)
	// r.GET("/users/:userId/submissions", getSubmissions)
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

// type Submission struct {
// 	Content string `json:"content"`
// 	Origin  string `json:"origin"`
// }

// func getSubmissions(c *gin.Context) {
// 	userId := c.Param("userId")

// 	sql := `
// 		select content, origin
// 		from submission
// 		where user_id = ?
// 	`

// 	rows, err := db.QueryContext(c, sql, userId)
// 	if err != nil {
// 		logger.Error("error fetching submissions: ", err)
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var submissions []Submission

// 	for rows.Next() {
// 		var submission Submission
// 		err := rows.Scan(&submission.Content, &submission.Origin)
// 		if err != nil {
// 			logger.Error("error scanning submission: ", err)
// 			c.AbortWithStatus(http.StatusInternalServerError)
// 			return
// 		}

// 		submissions = append(submissions, submission)
// 	}

// 	c.JSON(http.StatusOK, submissions)
// }

// type SubmissionPostRequestBody struct {
// 	Content string `json:"content"`
// }
