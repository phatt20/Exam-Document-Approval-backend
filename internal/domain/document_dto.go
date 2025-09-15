package domain

type CreateDocumentInput struct {
	DocumentName string `json:"document_name" validate:"required"`
}