package logger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)


func Test_Info(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{URL: &url.URL{Path: "url-path"}, Method: "req-method"}

	var b bytes.Buffer
	customLogger = zerolog.New(&b)
	SetLoggerContext(c)

	Info("Test message", " one two")
	expected := "{\"level\":\"info\",\"method\":\"req-method\",\"path\":\"url-path\",\"time\":\".*\",\"message\":\"Test message one two\"}"
	got := b.String()

	assert.Regexp(t, expected, got)
}

func Test_Error(t *testing.T) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{URL: &url.URL{Path: "url-path"}, Method: "req-method"}

	var b bytes.Buffer
	customLogger = zerolog.New(&b)
	SetLoggerContext(c)

	Error("error: ", fmt.Errorf("one error"))
	expected := "{\"level\":\"error\",\"method\":\"req-method\",\"path\":\"url-path\",\"error\":\"error: one error\",\"time\":\".*\"}"
	got := b.String()

	assert.Regexp(t, expected, got)
}