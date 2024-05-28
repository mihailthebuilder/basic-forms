package logger

import (
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var customLogger = zerolog.New(os.Stdout)

func SetLoggerContext(c *gin.Context) {
	customLogger = customLogger.With().Str("method", c.Request.Method).Str("path", c.Request.URL.Path).Logger()
	c.Next()
}

func Error(messages ...any) {
	clm := combineLoggedMessages(messages...)
	customLogger.Error().Err(errors.New(clm)).Timestamp().Msg("")
}

func Info(messages ...any) {
	clm := combineLoggedMessages(messages...)
	customLogger.Info().Timestamp().Msg(clm)
}

func combineLoggedMessages(messages ...any) string {
	return fmt.Sprint(messages...)
}
