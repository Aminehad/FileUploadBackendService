package repository

import (
	"errors"
	"upload-service/internal/app/upload-service/models"

	"gorm.io/gorm"
)

var (
	ErrRecordNotFound                = gorm.ErrRecordNotFound
	ErrDuplicatedKeyUniqueConstraint = errors.New("duplicate key value violates unique constraint")
)

type repository struct {
	db *gorm.DB
}

//go:generate mockery --name Repository
type Repository interface {
	BeginTx() (Repository, error)
	Commit() error
	Rollback() error
	AutoMigrate() error
	// TODO: methods for file CRUD operations
}

func New(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) AutoMigrate() error {
	if err := r.db.AutoMigrate(
		&models.File{},
	); err != nil {
		return err
	}
	return nil
}

func (r *repository) BeginTx() (Repository, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &repository{db: tx}, nil
}

func (r *repository) Commit() error {
	return r.db.Commit().Error
}

func (r *repository) Rollback() error {
	return r.db.Rollback().Error
}
