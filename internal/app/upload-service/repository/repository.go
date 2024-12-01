package repository

import (
	"errors"
	"strings"
	"upload-service/internal/app/upload-service/models"

	"github.com/google/uuid"
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
	SaveFileMetaData(file models.File) (*uuid.UUID, error)
	ListFiles(f FileFilters) ([]models.File, error)
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

func (r *repository) SaveFileMetaData(file models.File) (*uuid.UUID, error) {
	err := r.db.Create(&file).Error
	if err != nil {
		if strings.Contains(err.Error(), ErrDuplicatedKeyUniqueConstraint.Error()) {
			return nil, ErrDuplicatedKeyUniqueConstraint
		}
	}
	return &file.ID, nil
}

func (r *repository) ListFiles(f FileFilters) ([]models.File, error) {
	var files []models.File
	tx := r.db.Model(&models.File{}).Where(&models.File{
		ID: f.ID,
	})
	err := tx.Find(&files).Error
	if len(files) == 0 {
		return nil, ErrRecordNotFound
	}
	return files, err
}

type FileFilters struct {
	ID uuid.UUID
}
