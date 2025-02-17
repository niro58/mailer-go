package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

/*
{
    code : int,
    data : {
        ...
    },
    message :"OK"|"ERROR"
}
*/

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrUnauthorized   = errors.New("unauthorized")
)
func CreateReply(data interface{}, err error) gin.H {
	if err == nil {
		return gin.H{
			"error":   "OK",
			"data":    data,
			"message": "",
		}
	}else{
		return gin.H{
			"error":   "ERROR",
			"data":    data,
			"message": err.Error(),
		}
	}
}
