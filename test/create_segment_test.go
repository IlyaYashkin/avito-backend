package test

import (
	"avito-backend/internal/database"
	"avito-backend/internal/entity/segment"
	"avito-backend/utils"
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateSegment(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	router := gin.Default()
	router.POST("/create-segment", segment.CreateSegment)

	requestJson := []byte(`{"name": "AVITO_VOICE_100", "user_percentage": 33}`)
	req, err := http.NewRequest("POST", "/create-segment", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.JSONEq(t, `{"data":{"message":"Segment created","name":"AVITO_VOICE_100"},"status":"success"}`, w.Body.String())
}
