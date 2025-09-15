package usecase

import (
	"context"
	"errors"
	"testing"

	"approval-system/internal/domain"

	"github.com/stretchr/testify/assert"
)

// --- Mock Repository ---
type mockDocRepo struct{}

func (m *mockDocRepo) CreateDoc(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
	doc.ID = 1
	doc.Status = domain.DocumentStatusPending
	return doc, nil
}

func (m *mockDocRepo) FindAllDocs(ctx context.Context) ([]*domain.Document, error) {
	return []*domain.Document{
		{ID: 1, DocumentName: "Doc 1", Status: domain.DocumentStatusPending},
		{ID: 2, DocumentName: "Doc 2", Status: domain.DocumentStatusApproved},
	}, nil
}

func (m *mockDocRepo) FindDocByID(ctx context.Context, id int64) (*domain.Document, error) {
	if id == 0 {
		return nil, errors.New("not found")
	}
	return &domain.Document{ID: id, DocumentName: "Doc Test", Status: domain.DocumentStatusPending}, nil
}

func (m *mockDocRepo) UpdateDocsStatus(ctx context.Context, ids []int64, status domain.DocumentStatus, reason string) ([]*domain.Document, error) {
	if len(ids) == 0 {
		return nil, errors.New("no ids provided")
	}
	return []*domain.Document{
		{ID: ids[0], DocumentName: "Doc Test", Status: status},
	}, nil
}

// --- Setup ---
func setupUsecase() *docUsecase {
	return &docUsecase{docRepo: &mockDocRepo{}}
}

// --- Tests happynaja  ---
func TestCreateDoc(t *testing.T) {
	u := setupUsecase()

	input := &domain.CreateDocumentInput{DocumentName: "New Doc"}
	res, err := u.CreateDoc(context.Background(), input)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.ID)
	assert.Equal(t, "New Doc", res.DocumentName)
	assert.Equal(t, domain.DocumentStatusPending, res.Status)
}

func TestFindAllDocs(t *testing.T) {
	u := setupUsecase()

	res, err := u.FindAllDocs(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestFindDocByID(t *testing.T) {
	u := setupUsecase()

	res, err := u.FindDocByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.ID)

	_, err = u.FindDocByID(context.Background(), 0)
	assert.Error(t, err)
	assert.Equal(t, "document not found", err.Error())
}

func TestUpdateStatus(t *testing.T) {
	u := setupUsecase()

	input := &domain.UpdateStatusInput{
		DocumentIDs: []int64{1},
		Status:      "APPROVED",
		Reason:      "ok",
	}

	res, err := u.UpdateStatus(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, domain.DocumentStatusApproved, res[0].Status)

	input.Status = "INVALID"
	_, err = u.UpdateStatus(context.Background(), input)
	assert.Error(t, err)
	assert.Equal(t, "invalid status", err.Error())
}
