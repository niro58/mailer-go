package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrNoFile         = errors.New("no file uploaded")
	ErrServer         = errors.New("server error")
	ErrFile           = errors.New("file error")
)

func Respond(c *gin.Context, data interface{}, err error) {
	if err != nil {
		var statusCode int
		switch err {
		case ErrInvalidRequest:
			statusCode = http.StatusBadRequest
		case ErrUnauthorized:
			statusCode = http.StatusUnauthorized
		case ErrNoFile:
			statusCode = http.StatusBadRequest
		case ErrServer:
			statusCode = http.StatusInternalServerError
		case ErrFile:
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}
	if data != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"data":  data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"error": nil,
			"data":  "ok",
		})
	}

}
