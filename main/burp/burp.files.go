package burp

import (
	"net/http"
	"os"
	"path"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		responser.ResponseJson(c, http.StatusInternalServerError, "", "No file uploaded")
		return
	}

	// The file cannot be received.
	project, _ := c.GetPostForm("project")
	id, _ := c.GetPostForm("id")

	// Create folder for serve static files
	dir := path.Join("./files", project)
	_ = os.Mkdir(dir, 0755)

	fileName := path.Join(dir, id)

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, fileName); err != nil {
		responser.ResponseJson(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	// File saved successfully. Return proper result
	responser.ResponseJson(c, http.StatusOK, "Your file has been successfully uploaded.", "")
}
