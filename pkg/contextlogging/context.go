package contextlogging

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ContextLogger contains logrus entries used when logging
type ContextLogger struct {
	*logrus.Entry
	ID string
}

// New creates and returns a new ContextLogger
func New() ContextLogger {
	return NewWithID(uuid.New().String())
}

// NewWithID creates and returns a new ContextLogger with specified ID
func NewWithID(id string) ContextLogger {
	return ContextLogger{
		logrus.WithField("contextId", id),
		id,
	}
}

// GetGinLogContext - get log context from gin
func GetGinLogContext(c *gin.Context) (cc ContextLogger) {
	data, ok := c.Get("log")
	if !ok {
		log.Fatal("log context not set")
	}

	logger, ok := data.(ContextLogger)
	if !ok {
		log.Fatal("log was not logging.ContextLogger")
	}

	return logger
}
