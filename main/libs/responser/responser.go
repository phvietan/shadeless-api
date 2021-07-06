package responser

import "github.com/gin-gonic/gin"

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
