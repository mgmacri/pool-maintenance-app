package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/mgmacri/pool-maintenance-app/internal/domain"
)

type PostgresPoolRepository struct {
	DB *sql.DB
}

func NewPostgresPoolRepository(db *sql.DB) *PostgresPoolRepository {
	return &PostgresPoolRepository{DB: db}
}

func (r *PostgresPoolRepository) Create(p domain.Pool) error {
	// If CreatedAt is zero, let DB use default?
	// Easier to just always set it in Service or rely on Go time.
	// But here I'll just insert everything.
	_, err := r.DB.Exec("INSERT INTO pools (id, name, gallons, type, created_at) VALUES ($1, $2, $3, $4, $5)",
		p.ID, p.Name, p.Gallons, string(p.Type), p.CreatedAt)
	return err
}

func (r *PostgresPoolRepository) GetByID(id string) (domain.Pool, error) {
	row := r.DB.QueryRow("SELECT id, name, gallons, type, created_at FROM pools WHERE id = $1", id)
	var p domain.Pool
	var typeStr string
	err := row.Scan(&p.ID, &p.Name, &p.Gallons, &typeStr, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.Pool{}, ErrPoolNotFound
	} else if err != nil {
		return domain.Pool{}, err
	}
	p.Type = domain.PoolType(typeStr)
	return p, nil
}

func (r *PostgresPoolRepository) List() ([]domain.Pool, error) {
	rows, err := r.DB.Query("SELECT id, name, gallons, type, created_at FROM pools")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pools []domain.Pool
	for rows.Next() {
		var p domain.Pool
		var typeStr string
		if err := rows.Scan(&p.ID, &p.Name, &p.Gallons, &typeStr, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.Type = domain.PoolType(typeStr)
		pools = append(pools, p)
	}
	return pools, nil
}
