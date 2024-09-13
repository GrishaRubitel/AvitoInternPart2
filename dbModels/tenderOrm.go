package models

import (
	"time"
)

// TenderStatus ENUM
type TenderStatus string

const (
	Created   TenderStatus = "Created"
	Published TenderStatus = "Published"
	Closed    TenderStatus = "Closed"
)

// TenderServiceType ENUM
type TenderServiceType string

const (
	Construction TenderServiceType = "Construction"
	Delivery     TenderServiceType = "Delivery"
	Manufacture  TenderServiceType = "Manufacture"
)

// Tender table
type Tender struct {
	Id             string            `gorm:"type:varchar(100);primaryKey" json:"id"`               // tenderId
	Name           string            `gorm:"type:varchar(100);not null" json:"name"`               // tenderName
	Description    string            `gorm:"type:varchar(500);not null" json:"description"`        // tenderDescription
	ServiceType    TenderServiceType `gorm:"type:tender_service_type;not null" json:"serviceType"` // tenderServiceType (ENUM)
	Status         TenderStatus      `gorm:"type:tender_status;not null" json:"status"`            // tenderStatus (ENUM)
	OrganizationId string            `gorm:"type:varchar(100);not null" json:"organizationId"`     // organizationId
	Version        int               `gorm:"not null;default:1;check:version >= 1" json:"version"` // tenderVersion
	CreatedAt      time.Time         `gorm:"type:timestamp;not null" json:"createdAt"`             // createdAt
}
