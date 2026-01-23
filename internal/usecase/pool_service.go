package usecase

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mgmacri/pool-maintenance-app/internal/domain"
	"github.com/mgmacri/pool-maintenance-app/internal/repository"
)

// PoolUseCase defines the core business logic.
type PoolUseCase struct {
	repo repository.PoolRepository
}

func NewPoolUseCase(repo repository.PoolRepository) *PoolUseCase {
	return &PoolUseCase{repo: repo}
}

func (uc *PoolUseCase) RegisterPool(name string, gallons int, pType domain.PoolType) (domain.Pool, error) {
	pool := domain.Pool{
		ID:      uuid.New().String(),
		Name:    name,
		Gallons: gallons,
		Type:    pType,
	}
	if err := uc.repo.Create(pool); err != nil {
		return domain.Pool{}, err
	}
	return pool, nil
}

func (uc *PoolUseCase) GetPool(id string) (domain.Pool, error) {
	return uc.repo.GetByID(id)
}

func (uc *PoolUseCase) AnalyzeChemistry(readings domain.ChemicalReading) (domain.JobReport, error) {
	// Verify pool exists
	pool, err := uc.repo.GetByID(readings.PoolID)
	if err != nil {
		return domain.JobReport{}, err
	}

	report := domain.JobReport{
		PoolID:     pool.ID,
		IsBalanced: true,
		Dosage:     []string{},
	}

	// rudimentary dosing logic (Command Pattern placeholder)
	if readings.PH < 7.2 {
		report.IsBalanced = false
		report.Dosage = append(report.Dosage, "Add Soda Ash to raise pH")
	} else if readings.PH > 7.8 {
		report.IsBalanced = false
		report.Dosage = append(report.Dosage, "Add Muriatic Acid to lower pH")
	}

	if readings.FreeChlorine < 1.0 {
		report.IsBalanced = false
		amount := float64(pool.Gallons) * 0.00013 // rough calc for 1ppm increase using liquid chlorine
		report.Dosage = append(report.Dosage, fmt.Sprintf("Add %.2f gal Liquid Chlorine", amount))
	}

	return report, nil
}
