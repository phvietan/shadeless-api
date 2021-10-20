package burp

import (
	"github.com/gin-gonic/gin"
)

func Routes(route *gin.Engine) {
	burp := route.Group("/burp")
	{
		burp.POST("/packets", postPackets)
		burp.POST("/files", uploadFile)
	}
}
