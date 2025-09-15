package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"approval-system/internal/domain"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockDocUsecase struct{}

func (m *mockDocUsecase) CreateDoc(ctx context.Context, input *domain.CreateDocumentInput) (*domain.Document, error) {
	return &domain.Document{ID: 1, DocumentName: input.DocumentName, Status: domain.DocumentStatusPending}, nil
}

func (m *mockDocUsecase) FindAllDocs(ctx context.Context) ([]*domain.Document, error) {
	return []*domain.Document{
		{ID: 1, DocumentName: "Doc 1", Status: domain.DocumentStatusPending},
		{ID: 2, DocumentName: "Doc 2", Status: domain.DocumentStatusApproved},
	}, nil
}

func (m *mockDocUsecase) FindDocByID(ctx context.Context, id int64) (*domain.Document, error) {
	if id == 0 {
		return nil, errors.New("not found")
	}
	return &domain.Document{ID: id, DocumentName: "Doc Test", Status: domain.DocumentStatusPending}, nil
}

func (m *mockDocUsecase) UpdateStatus(ctx context.Context, input *domain.UpdateStatusInput) ([]*domain.Document, error) {
	return []*domain.Document{
		{ID: input.DocumentIDs[0], DocumentName: "Doc Test", Status: domain.DocumentStatus(input.Status)},
	}, nil
}

// --- Tests --- เริ่ม imp
func setupHandler() *docHttpHandler {
	return &docHttpHandler{
		docUsecase: &mockDocUsecase{},
	}
}

func TestCreateDoc(t *testing.T) {
	e := echo.New()
	handler := setupHandler()

	body := map[string]interface{}{"document_name": "Test Doc"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/doc/create", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateDoc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var res domain.Document
	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, "Test Doc", res.DocumentName)
}

func TestFindAllDocs(t *testing.T) {
	e := echo.New()
	handler := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/doc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.FindAllDocs(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res []*domain.Document
	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Len(t, res, 2)
}

func TestFindDocByID(t *testing.T) {
	e := echo.New()
	handler := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/doc/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := handler.FindDocByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res domain.Document
	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, int64(1), res.ID)
}

func TestUpdateStaus(t *testing.T) {
	e := echo.New()
	handler := setupHandler()

	body := map[string]interface{}{
		"document_ids": []int64{1},
		"status":       "APPROVED",
		"reason":       "ok",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/doc/update-staus", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.UpdateStaus(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res []domain.Document
	json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Equal(t, "APPROVED", string(res[0].Status))
}

// เคสที่ต้อง error
func TestCreateDoc_Error(t *testing.T) {
	e := echo.New()
	handler := &docHttpHandler{
		docUsecase: &mockDocUsecase{},
	}

	req := httptest.NewRequest(http.MethodPost, "/doc/create", bytes.NewReader([]byte("{}")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.CreateDoc(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFindDocByID_NotFound(t *testing.T) {
	e := echo.New()
	handler := setupHandler()

	req := httptest.NewRequest(http.MethodGet, "/doc/0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("0")

	err := handler.FindDocByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
