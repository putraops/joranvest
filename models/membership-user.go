package models

import (
	"database/sql"
)

type MembershipUser struct {
	Id          string       `gorm:"type:varchar(50);primary_key" json:"id"`
	IsActive    bool         `gorm:"type:bool;default:1" json:"is_active"`
	IsLocked    bool         `gorm:"type:bool" json:"is_locked"`
	IsDefault   bool         `gorm:"type:bool" json:"is_default"`
	CreatedAt   sql.NullTime `gorm:"type:timestamp" json:"created_at"`
	CreatedBy   string       `gorm:"type:varchar(50)" json:"created_by"`
	UpdatedAt   sql.NullTime `gorm:"type:timestamp" json:"updated_at"`
	UpdatedBy   string       `gorm:"type:varchar(50)" json:"updated_by"`
	SubmittedAt sql.NullTime `gorm:"type:timestamp" json:"submitted_at"`
	SubmittedBy string       `gorm:"type:varchar(50)" json:"submitted_by"`
	ApprovedAt  sql.NullTime `gorm:"type:timestamp;default:null" json:"approved_at"`
	ApprovedBy  string       `gorm:"type:varchar(50)" json:"approved_by"`
	OwnerId     string       `gorm:"type:varchar(50)" json:"owner_id"`
	EntityId    string       `gorm:"type:varchar(50);null" json:"entity_id"`

	MembershipId      string       `gorm:"type:varchar(50);not null" json:"membership_id"`
	ApplicationUserId string       `gorm:"type:varchar(50);not null" json:"application_user_id"`
	PaymentId         string       `gorm:"type:varchar(50);" json:"payment_id"`
	StartedDate       sql.NullTime `gorm:"type:timestamp" json:"started_date"`
	ExpiredDate       sql.NullTime `gorm:"type:timestamp" json:"expired_date"`

	Membership      Membership      `gorm:"foreignkey:MembershipId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"membership"`
	Payment         Payment         `gorm:"foreignkey:PaymentId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"payment"`
	ApplicationUser ApplicationUser `gorm:"foreignkey:ApplicationUserId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"application_user"`
}

func (MembershipUser) TableName() string {
	return "membership_user"
}
