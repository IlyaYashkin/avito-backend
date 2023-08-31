package test

import (
	"avito-backend/internal/database"
	"avito-backend/internal/entity/segment"
	"avito-backend/internal/entity/usersegment"
	"avito-backend/internal/entity/usersegmentlog"
	"avito-backend/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddUserSegments(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	router := gin.Default()
	router.POST("/update-user-segments", usersegment.UpdateUserSegments)

	userId := int32(1000)
	numSegments := 50

	var segmentStrs []string
	for i := 1; i < numSegments; i++ {
		segmentStr := fmt.Sprintf("AVITO_VOICE_%d", 10*i)
		segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentStr})
		segmentStrs = append(segmentStrs, segmentStr)
	}

	requestJson, err := json.Marshal(gin.H{
		"user_id":      userId,
		"add_segments": segmentStrs,
	})
	req, err := http.NewRequest("POST", "/update-user-segments", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	userSegmentRepo := usersegment.NewUserSegmentRepository(db)
	resultSegments, err := userSegmentRepo.GetSegmentsNamesByUserId(userId)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, segmentStrs, resultSegments)
}

func TestDeleteUserSegments(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	router := gin.Default()
	router.POST("/update-user-segments", usersegment.UpdateUserSegments)

	userId := int32(1000)
	numSegments := 100

	var segmentStrs []any
	for i := 1; i < numSegments; i++ {
		segmentStr := fmt.Sprintf("AVITO_VOICE_%d", 10*i)
		segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentStr})
		segmentStrs = append(segmentStrs, segmentStr)
	}
	_, err := usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: userId, AddSegments: segmentStrs})
	if err != nil {
		t.Fatal(err)
	}

	requestJson, err := json.Marshal(gin.H{
		"user_id":         1000,
		"delete_segments": segmentStrs,
	})
	req, err := http.NewRequest("POST", "/update-user-segments", bytes.NewBuffer(requestJson))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	userSegmentRepo := usersegment.NewUserSegmentRepository(db)
	userSegments, err := userSegmentRepo.GetByUserId(1000)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, len(userSegments))
}

func TestGetUserSegments(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	router := gin.Default()
	router.GET("/get-user-segments/:user_id", usersegment.GetUserSegments)

	userId := int32(1000)
	numSegments := 100

	var segmentStrs []any
	for i := 1; i < numSegments; i++ {
		segmentStr := fmt.Sprintf("AVITO_VOICE_%d", 10*i)
		segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentStr})
		segmentStrs = append(segmentStrs, segmentStr)
	}
	_, err := usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: userId, AddSegments: segmentStrs})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("/get-user-segments/%d", userId), nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expectedArr, err := json.Marshal(segmentStrs)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(`{"data":{"segments":%s,"user":%d},"status":"success"}`, string(expectedArr), userId), w.Body.String())
}

// OPTIONAL 1

func TestUserSegmentLog(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	router := gin.Default()
	router.GET("/get-user-segment-log", usersegmentlog.GetUserSegmentLog)

	userId := int32(1000)
	numSegments := 50

	var addSegments []any
	var deleteSegments []string
	for i := 1; i < numSegments; i++ {
		segmentStr := fmt.Sprintf("AVITO_VOICE_%d", 10*i)
		segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentStr})
		addSegments = append(addSegments, segmentStr)
		deleteSegments = append(deleteSegments, segmentStr)
	}
	_, err := usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: userId, AddSegments: addSegments})
	if err != nil {
		t.Fatal(err)
	}
	_, err = usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: userId, DeleteSegments: deleteSegments})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/get-user-segment-log", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	mockRequest := usersegmentlog.RequestGetUserSegmentLog{}

	userSegmentLogRepo := usersegmentlog.NewUserSegmentLogRepository(db)
	logRows, err := userSegmentLogRepo.Get(mockRequest.UserId, mockRequest.Date)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(addSegments)+len(deleteSegments), len(logRows))
}

func TestUserSegmentLogUserId(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	router := gin.Default()
	router.GET("/get-user-segment-log", usersegmentlog.GetUserSegmentLog)

	userId := int32(1000)
	numSegments := 50

	var addSegments []any
	for i := 1; i < numSegments; i++ {
		segmentStr := fmt.Sprintf("AVITO_VOICE_%d", 10*i)
		segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentStr})
		addSegments = append(addSegments, segmentStr)
	}
	_, err := usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: userId, AddSegments: addSegments})
	if err != nil {
		t.Fatal(err)
	}
	_, err = usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: 1001, AddSegments: addSegments})
	if err != nil {
		t.Fatal(err)
	}
	_, err = usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: 1002, AddSegments: addSegments})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/get-user-segment-log", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	mockRequest := usersegmentlog.RequestGetUserSegmentLog{UserId: 1000}

	userSegmentLogRepo := usersegmentlog.NewUserSegmentLogRepository(db)
	logRows, err := userSegmentLogRepo.Get(mockRequest.UserId, mockRequest.Date)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(addSegments), len(logRows))
}

func TestUserSegmentLogDate(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	router := gin.Default()
	router.GET("/get-user-segment-log", usersegmentlog.GetUserSegmentLog)

	userId := int32(1000)
	numSegments := 50

	var addSegments []any
	for i := 1; i < numSegments; i++ {
		segmentStr := fmt.Sprintf("AVITO_VOICE_%d", 10*i)
		segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentStr})
		addSegments = append(addSegments, segmentStr)
	}
	_, err := usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: userId, AddSegments: addSegments})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/get-user-segment-log", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	if err != nil {
		t.Fatal(err)
	}
	date := time.Now()
	date = date.AddDate(0, -1, 0)
	mockRequest := usersegmentlog.RequestGetUserSegmentLog{Date: date}

	userSegmentLogRepo := usersegmentlog.NewUserSegmentLogRepository(db)
	logRows, err := userSegmentLogRepo.Get(mockRequest.UserId, mockRequest.Date)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 0, len(logRows))
}

// OPTIONAL 2

func TestUserSegmentTtl(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	router := gin.Default()
	router.GET("/get-user-segments/:user_id", usersegment.GetUserSegments)

	userId := int32(1000)
	numSegments := 50

	var addTtlSegments []any
	for i := 1; i < numSegments; i++ {
		segmentStr := fmt.Sprintf("AVITO_VOICE_%d", 10*i)
		segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentStr})
		date := time.Now()
		date = date.AddDate(0, 0, -1)
		addTtlSegments = append(addTtlSegments, map[string]interface{}{"segment": segmentStr, "ttl": date.Format(usersegment.TTL_TIME_FORMAT)})
	}
	var addSegments []any
	for i := numSegments; i < numSegments*2; i++ {
		segmentStr := fmt.Sprintf("AVITO_VOICE_%d", 10*i)
		segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentStr})
		addSegments = append(addSegments, segmentStr)
	}

	_, err := usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: userId, AddSegments: addTtlSegments})
	if err != nil {
		t.Fatal(err)
	}
	_, err = usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: userId, AddSegments: addSegments})
	if err != nil {
		t.Fatal(err)
	}
	usersegment.ClearExpiredTtlUserSegments(db)

	req, err := http.NewRequest("GET", fmt.Sprintf("/get-user-segments/%d", userId), nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	expectedArr, err := json.Marshal(addSegments)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf(`{"data":{"segments":%s,"user":%d},"status":"success"}`, string(expectedArr), userId), w.Body.String())
}

// OPTIONAL 3

func TestSegmentPercentage(t *testing.T) {
	db := database.Open()
	defer db.Close()

	utils.ClearDB()

	segmentName := "AVITO_VOICE_30"
	segmentPercentage := 10
	numUsers := 10

	err := segment.CreateSegmentService(segment.RequestUpdateSegment{Name: segmentName, UserPercentage: float32(segmentPercentage)})
	if err != nil {
		t.Fatal(err)
	}

	for i := 1000; i < 1000+numUsers; i++ {
		_, err := usersegment.UpdateUserSegmentsService(usersegment.RequestUpdateUserSegments{UserId: int32(i)})
		if err != nil {
			t.Fatal(err)
		}
	}

	userSegmentRepo := usersegment.NewUserSegmentRepository(db)

	userSegments, err := userSegmentRepo.Get()
	if err != nil {
		t.Fatal(err)
	}

	expectedRowsNum := float64(segmentPercentage) / float64(100) * float64(numUsers)
	assertion := int(math.Abs(expectedRowsNum-float64(len(userSegments)))) <= 1

	assert.Equal(t, true, assertion)
}
