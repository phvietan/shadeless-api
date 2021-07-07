package responser

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type responseBody struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
}

func ResponseJson(c *gin.Context, status int, data interface{}, error string) {
	response := responseBody{
		StatusCode: status,
		Data:       data,
		Error:      error,
	}
	c.JSON(status, response)
}

func ResponseOk(c *gin.Context, data interface{}) {
	response := responseBody{
		StatusCode: http.StatusOK,
		Data:       data,
		Error:      "",
	}
	c.JSON(http.StatusOK, response)
}

func ResponseError(c *gin.Context, err error) {
	response := responseBody{
		StatusCode: http.StatusInternalServerError,
		Data:       "",
		Error:      err.Error(),
	}
	c.JSON(http.StatusInternalServerError, response)
}
