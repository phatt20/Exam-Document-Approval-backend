package repository

import (
	"approval-system/internal/domain"
	"approval-system/pkg/database"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type DocRepositoryInterface interface {
	CreateDoc(ctx context.Context, req *domain.Document) (*domain.Document, error)
	FindAllDocs(ctx context.Context) ([]*domain.Document, error)
	FindDocByID(ctx context.Context, id int64) (*domain.Document, error)
	UpdateDocsStatus(ctx context.Context, ids []int64, status domain.DocumentStatus, reason string) ([]*domain.Document, error)
}

type docRepository struct {
	db database.DatabasesPostgres
}

func NewDocRepository(db database.DatabasesPostgres) DocRepositoryInterface {
	return &docRepository{db}
}

func (r *docRepository) CreateDoc(ctx context.Context, req *domain.Document) (*domain.Document, error) {
	db := r.db.Connect()

	req.Status = domain.DocumentStatusPending

	if err := db.WithContext(ctx).Create(req).Error; err != nil {
		return nil, err
	}

	return req, nil
}

func (r *docRepository) FindAllDocs(ctx context.Context) ([]*domain.Document, error) {
	db := r.db.Connect()
	var docs []*domain.Document

	if err := db.WithContext(ctx).Order("created_at DESC").Find(&docs).Error; err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *docRepository) FindDocByID(ctx context.Context, id int64) (*domain.Document, error) {
	db := r.db.Connect()
	var doc domain.Document

	if err := db.WithContext(ctx).First(&doc, id).Error; err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r *docRepository) UpdateDocsStatus(ctx context.Context, ids []int64, status domain.DocumentStatus, reason string) ([]*domain.Document, error) {
	var docs []*domain.Document

	if err := r.db.Connect().WithContext(ctx).Where("id IN ?", ids).Find(&docs).Error; err != nil {
		return nil, err
	}

	for _, d := range docs {
		if d.Status != domain.DocumentStatusPending {
			return nil, fmt.Errorf("document already processed: %d", d.ID)
		}
	}

	err := r.db.Connect().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(&domain.Document{}).
			Where("id IN ?", ids).
			Updates(map[string]interface{}{
				"status":     status,
				"reason":     reason,
				"updated_at": gorm.Expr("NOW()"),
			}).Error
	})
	if err != nil {
		return nil, err
	}

	return docs, nil
}
