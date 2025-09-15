package domain

import "time"

type DocumentStatus string

const (
	DocumentStatusPending  DocumentStatus = "PENDING" //สถานะเริ่มต้น
	DocumentStatusApproved DocumentStatus = "APPROVED"
	DocumentStatusRejected DocumentStatus = "REJECTED"
)

type Document struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentName string         `gorm:"not null" json:"document_name"`
	Status       DocumentStatus `gorm:"type:varchar(16);not null;default:'PENDING'" json:"status"`
	Reason       string         `gorm:"type:text;default:null" json:"reason"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Document) TableName() string {
	return "documents"
}

type UpdateStatusInput struct {
	DocumentIDs []int64 `json:"document_ids" validate:"required,min=1"`
	Reason      string  `json:"reason" validate:"required"`
	Status      string  `json:"status" validate:"required,oneof=PENDING APPROVED REJECTED"`
}
