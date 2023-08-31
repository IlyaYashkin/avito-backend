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

func TestDeleteSegment(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	segment.CreateSegmentService(segment.RequestUpdateSegment{Name: "AVITO_VOICE_100"})

	router := gin.Default()
	router.POST("/delete-segment", segment.DeleteSegment)

	requestJson := []byte(`{"name": "AVITO_VOICE_100"}`)
	req, err := http.NewRequest("POST", "/delete-segment", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"data":{"message":"Segment deleted","name":"AVITO_VOICE_100"},"status":"success"}`, w.Body.String())
}
