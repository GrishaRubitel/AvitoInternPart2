package models

import (
	"time"
)

// BidStatus ENUM

type BidStatus string

const (
	BidStatusCreated   BidStatus = "Created"
	BidStatusPublished BidStatus = "Published"
	BidStatusCanceled  BidStatus = "Canceled"
	BidStatusApproved  BidStatus = "Approved"
	BidStatusRejected  BidStatus = "Rejected"
)

// BidAuthorType ENUM
type BidAuthorType string

const (
	BidAuthorTypeOrganization BidAuthorType = "Organization"
	BidAuthorTypeUser         BidAuthorType = "User"
)

// Bid table
type Bid struct {
	Id          string        `gorm:"type:varchar(100);primaryKey" json:"id"`
	Name        string        `gorm:"type:varchar(100);not null" json:"name"`
	Description string        `gorm:"type:varchar(500);not null" json:"description"`
	Status      BidStatus     `gorm:"type:bid_status;not null" json:"status"`          // ENUM type
	TenderId    string        `gorm:"type:varchar(100);index" json:"tenderId"`         // Foreign Key
	AuthorType  BidAuthorType `gorm:"type:bid_author_type;not null" json:"authorType"` // ENUM type
	AuthorId    string        `gorm:"type:varchar(100);not null" json:"authorId"`
	Version     int           `gorm:"default:1;not null" json:"version"`
	CreatedAt   time.Time     `gorm:"type:timestamp;not null" json:"createdAt"`
}

// BidDecision ENUM
type BidDecision string

const (
	BidDecisionApproved BidDecision = "Approved"
	BidDecisionRejected BidDecision = "Rejected"
)

// BidReview table
type BidReview struct {
	Id          string      `gorm:"primaryKey;type:varchar(100)"` //bidReviewId
	Description string      `gorm:"type:varchar(1000);not null"`  //bidReviewDecision
	CreatedAt   time.Time   `gorm:"type:timestamp;not null"`      //createdAt
	BidId       string      `gorm:"type:varchar(100);index"`      // Foreign Key
	Decision    BidDecision `gorm:"type:bid_decision;not null"`   // ENUM type - This attribute was added by myself, because it appeared in openapi.yaml, but not in any table
}
