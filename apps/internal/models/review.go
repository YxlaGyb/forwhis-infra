package models

type Review struct {
	ReviewId   int    `gorm:"primaryKey;autoIncrement"`
	TargetType string `gorm:"size:50;not null"`
	TargetId   int    `gorm:"not null"`
	ReviewerId int
	Status     int    `gorm:"default:0"`
	Reason     string `gorm:"type:text"`
	Data       string `gorm:"type:text"`
	CreatedAt  int    `gorm:"default:0"`
	UpdatedAt  int    `gorm:"default:0"`
}

const (
	ReviewStatusPending  = 0
	ReviewStatusApproved = 1
	ReviewStatusRejected = 2
)

const (
	TargetTypeUserProfile = "user_profile"
)

func (Review) TableName() string { return "reviews" }
