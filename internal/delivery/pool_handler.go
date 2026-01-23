package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mgmacri/pool-maintenance-app/internal/domain"
	"github.com/mgmacri/pool-maintenance-app/internal/usecase"
	"go.uber.org/zap"
)

// PoolHandler handles HTTP requests for pool management.
type PoolHandler struct {
	logger  *zap.Logger
	service *usecase.PoolUseCase
}

// NewPoolHandler creates a new PoolHandler.
func NewPoolHandler(logger *zap.Logger, service *usecase.PoolUseCase) *PoolHandler {
	return &PoolHandler{
		logger:  logger,
		service: service,
	}
}

type CreatePoolRequest struct {
	Name    string          `json:"name" binding:"required"`
	Gallons int             `json:"gallons" binding:"required,min=100"`
	Type    domain.PoolType `json:"type" binding:"required,oneof=chlorine salt"`
}

// CreatePool registers a new pool.
// @Summary Register a new pool
// @Description Creates a new pool entity in the system
// @Tags pools
// @Accept json
// @Produce json
// @Param request body CreatePoolRequest true "Pool details"
// @Success 201 {object} domain.Pool
// @Failure 400 {object} map[string]string
// @Router /api/v1/pools [post]
func (h *PoolHandler) CreatePool(c *gin.Context) {
	var req CreatePoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid pool creation request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pool, err := h.service.RegisterPool(req.Name, req.Gallons, req.Type)
	if err != nil {
		h.logger.Error("failed to create pool", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register pool"})
		return
	}

	h.logger.Info("pool created", zap.String("pool_id", pool.ID))
	c.JSON(http.StatusCreated, pool)
}

// GetPool retrieves a pool by ID.
// @Summary Get a pool
// @Description Retrieves pool details by ID
// @Tags pools
// @Produce json
// @Param id path string true "Pool ID"
// @Success 200 {object} domain.Pool
// @Failure 404 {object} map[string]string
// @Router /api/v1/pools/{id} [get]
func (h *PoolHandler) GetPool(c *gin.Context) {
	id := c.Param("id")
	pool, err := h.service.GetPool(id)
	if err != nil {
		h.logger.Warn("pool not found", zap.String("pool_id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "pool not found"})
		return
	}
	c.JSON(http.StatusOK, pool)
}

// AnalyzeReadings submits chemical readings and gets a treatment plan.
// @Summary Analyze water chemistry
// @Description Submits readings (pH, Cl) and returns dosage instructions based on pool volume
// @Tags chemistry
// @Accept json
// @Produce json
// @Param request body domain.ChemicalReading true "Water readings"
// @Success 200 {object} domain.JobReport
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/chemistry/analyze [post]
func (h *PoolHandler) AnalyzeChemistry(c *gin.Context) {
	var req domain.ChemicalReading
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report, err := h.service.AnalyzeChemistry(req)
	if err != nil {
		h.logger.Warn("chemistry analysis failed", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("chemistry analyzed",
		zap.String("pool_id", req.PoolID),
		zap.Bool("balanced", report.IsBalanced),
	)
	c.JSON(http.StatusOK, report)
}
