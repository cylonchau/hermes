package query

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse handles successful API requests
func SuccessResponse(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusOK, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// APIResponse is a generic API response handler
func APIResponse(ctx *gin.Context, err error, data interface{}) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(http.StatusOK, Response{
		Code: returnCode,
		Msg:  message,
		Data: data,
	})
}

// ErrorResponse handles API errors with specific status codes
func ErrorResponse(ctx *gin.Context, statusCode int, err error) {
	returnCode, message := DecodeErr(err)
	ctx.JSON(statusCode, Response{
		Code: returnCode,
		Msg:  message,
	})
}

// InternalError handles 500 Internal Server Error
func InternalError(ctx *gin.Context, err error) {
	ErrorResponse(ctx, http.StatusInternalServerError, err)
}

// BadRequest handles 400 Bad Request
func BadRequest(ctx *gin.Context, err error) {
	ErrorResponse(ctx, http.StatusBadRequest, err)
}

// NotFound handles 404 Not Found
func NotFound(ctx *gin.Context, err error) {
	ErrorResponse(ctx, http.StatusNotFound, err)
}
