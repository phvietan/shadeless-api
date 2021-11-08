package burp

import (
	"errors"
	"fmt"
	"os"
	"path"
	"shadeless-api/main/config"
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/database/schema"
	"shadeless-api/main/libs/responser"

	"github.com/gin-gonic/gin"
)

func uploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error 1", err)
		responser.ResponseError(c, errors.New("No file uploaded"))
		return
	}

	// The file cannot be received.
	project, _ := c.GetPostForm("project")
	id, _ := c.GetPostForm("id")

	// Create folder for serve static files
	fileDir := config.GetInstance().GetFileDir()
	dir := path.Join(fileDir, project)
	_ = os.Mkdir(dir, 0755)

	fileName := path.Join(dir, id)

	// The file is received, so let's save it
	if err := c.SaveUploadedFile(file, fileName); err != nil {
		fmt.Println("Error 2", err)
		responser.ResponseError(c, err)
		return
	}

	var fileDatabase database.IFileDatabase = new(database.FileDatabase).Init()
	if fileInDb := fileDatabase.GetFileByProjectAndId(project, id); fileInDb == nil {
		newFileDB := schema.NewFile(project, id)
		if err := fileDatabase.Insert(newFileDB); err != nil {
			fmt.Println("Error 3:", err)
			responser.ResponseError(c, err)
			return
		}
	}
	responser.ResponseOk(c, "Your file has been successfully uploaded.")
}
