package projects

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shadeless-api/main/libs/database"
	"shadeless-api/main/libs/database/schema"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func spawnApp() *gin.Engine {
	router := gin.Default()
	Routes(router)
	return router
}

var router = spawnApp()

func TestGetProjectByName(t *testing.T) {
	type projectResponse struct {
		StatusCode int            `json:"statusCode"`
		Data       schema.Project `json:"data"`
		Error      string         `json:"error"`
	}

	var projectData database.IProjectDatabase = new(database.ProjectDatabase).Init()

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		projectName := "test"
		req, _ := http.NewRequest("GET", "/projects/"+projectName, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		var resp projectResponse

		err := json.Unmarshal([]byte(w.Body.String()), &resp)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.StatusCode, 404)
		assert.Equal(t, resp.Error, "Not found project with that name")

		objID, err := primitive.ObjectIDFromHex("000000000000000000000000")
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.Data.ID, primitive.ObjectID(objID))

		newProject := schema.NewProject()
		newProject.Name = projectName
		projectData.Insert(newProject)

		w2 := httptest.NewRecorder()
		req, _ = http.NewRequest("GET", "/projects/"+projectName, nil)
		router.ServeHTTP(w2, req)

		err = json.Unmarshal([]byte(w2.Body.String()), &resp)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp.StatusCode, 200)
		assert.Equal(t, resp.Error, "")
		assert.NotEqual(t, resp.Data.ID, primitive.ObjectID(objID))

		projectData.ClearCollection()
	}
}
