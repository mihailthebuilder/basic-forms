package main

import (
	"basic-forms/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setUpRoutes(r *gin.Engine) {
	r.GET("/users/:userId/submissions", getSubmissions)
	r.POST("/users/:userId/submit", postSubmission)
}

type Submission struct {
	Content   string `json:"content"`
	OriginUrl string `json:"originUrl"`
}

func getSubmissions(c *gin.Context) {
	userId := c.Param("userId")

	sql := `
		select content, origin_url
		from submission
		where user_id = ?
	`

	rows, err := db.QueryContext(c, sql, userId)
	if err != nil {
		logger.Error("error fetching submissions: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var submissions []Submission

	for rows.Next() {
		var submission Submission
		err := rows.Scan(&submission.Content, &submission.OriginUrl)
		if err != nil {
			logger.Error("error scanning submission: ", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		submissions = append(submissions, submission)
	}

	c.JSON(http.StatusOK, submissions)
}

type SubmissionPostRequestBody struct {
	Content string `json:"content"`
}

func postSubmission(c *gin.Context) {
	userId := c.Param("userId")
	host := c.Request.Header.Get("Origin")

	var body SubmissionPostRequestBody

	err := c.ShouldBindBodyWithJSON(&body)
	if err != nil {
		logger.Error("error parsing body: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	sql := `
		insert into submission (user_id, content, origin_url)
		values (?, ?, ?)
	`

	result, err := db.Exec(sql, userId, body.Content, host)

	if err != nil {
		logger.Error("error registering submission: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("error fetching rows affected: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if rows != 1 {
		logger.Error("rows affected error; expected 1, got: ", rows)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}
