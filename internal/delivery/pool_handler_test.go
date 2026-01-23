package delivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mgmacri/pool-maintenance-app/internal/delivery"
	"github.com/mgmacri/pool-maintenance-app/internal/domain"
	"github.com/mgmacri/pool-maintenance-app/internal/repository"
	"github.com/mgmacri/pool-maintenance-app/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupPoolRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	logger := zap.NewNop() // No logs during test
	repo := repository.NewInMemoryPoolRepo()
	svc := usecase.NewPoolUseCase(repo)
	handler := delivery.NewPoolHandler(logger, svc)

	r := gin.New()
	r.POST("/api/v1/pools", handler.CreatePool)
	r.POST("/api/v1/chemistry/analyze", handler.AnalyzeChemistry)
	return r
}

func TestCreatePool(t *testing.T) {
	r := setupPoolRouter()

	reqBody := delivery.CreatePoolRequest{
		Name:    "Test Pool",
		Gallons: 15000,
		Type:    domain.PoolTypeSalt,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/pools", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var pool domain.Pool
	err := json.Unmarshal(w.Body.Bytes(), &pool)
	assert.NoError(t, err)
	assert.NotEmpty(t, pool.ID)
	assert.Equal(t, "Test Pool", pool.Name)
}

func TestAnalyzeChemistry_BadPH(t *testing.T) {
	r := setupPoolRouter()

	// 1. Create a pool first (needed for lookups)
	// For integration test simplicity, we bypass HTTP and hit repo directly if we had access,
	// but here we just hit the endpoint.
	createBody := `{"name":"Dirty Pool","gallons":10000,"type":"chlorine"}`
	wCreate := httptest.NewRecorder()
	reqCreate, _ := http.NewRequest("POST", "/api/v1/pools", bytes.NewBufferString(createBody))
	r.ServeHTTP(wCreate, reqCreate)

	var pool domain.Pool
	_ = json.Unmarshal(wCreate.Body.Bytes(), &pool)

	// 2. Submit bad readings
	reading := domain.ChemicalReading{
		PoolID: pool.ID,
		PH:     6.8, // Low pH
	}
	body, _ := json.Marshal(reading)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/chemistry/analyze", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var report domain.JobReport
	err := json.Unmarshal(w.Body.Bytes(), &report)
	assert.NoError(t, err)
	assert.False(t, report.IsBalanced)
	assert.Contains(t, report.Dosage[0], "Soda Ash")
}
