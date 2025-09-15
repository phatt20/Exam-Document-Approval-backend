package usecase

import (
	docRepository "approval-system/internal/Repository"
	"approval-system/internal/domain"
	"context"
	"errors"
)

type DocUsecaseInterface interface {
	CreateDoc(ctx context.Context, req *domain.CreateDocumentInput) (*domain.Document, error)
	FindAllDocs(ctx context.Context) ([]*domain.Document, error)
	FindDocByID(ctx context.Context, id int64) (*domain.Document, error)
	UpdateStatus(ctx context.Context, req *domain.UpdateStatusInput) ([]*domain.Document, error)
}

type docUsecase struct {
	docRepo docRepository.DocRepositoryInterface
}

func NewDocUsecase(docRepo docRepository.DocRepositoryInterface) DocUsecaseInterface {
	return &docUsecase{docRepo}
}

func (u *docUsecase) CreateDoc(ctx context.Context, req *domain.CreateDocumentInput) (*domain.Document, error) {

	doc := &domain.Document{
		DocumentName: req.DocumentName,
	}
	return u.docRepo.CreateDoc(ctx, doc)
}

func (u *docUsecase) FindAllDocs(ctx context.Context) ([]*domain.Document, error) {
	return u.docRepo.FindAllDocs(ctx)
}

func (u *docUsecase) FindDocByID(ctx context.Context, id int64) (*domain.Document, error) {
	doc, err := u.docRepo.FindDocByID(ctx, id)
	if err != nil {
		return nil, errors.New("document not found")
	}
	return doc, nil
}

func (u *docUsecase) UpdateStatus(ctx context.Context, req *domain.UpdateStatusInput) ([]*domain.Document, error) {
	newStatus := domain.DocumentStatus(req.Status)
	if newStatus != domain.DocumentStatusApproved && newStatus != domain.DocumentStatusRejected {
		return nil, errors.New("invalid status")
	}

	return u.docRepo.UpdateDocsStatus(ctx, req.DocumentIDs, newStatus, req.Reason)
}
