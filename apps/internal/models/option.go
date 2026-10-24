package models

type Option struct {
	OptionId    int    `gorm:"primaryKey;autoIncrement"`
	OptionKey   string `gorm:"size:150;not null;uniqueIndex"`
	OptionValue string `gorm:"type:text"`
	Autoload    int    `gorm:"default:0"`
	CreatedAt   int    `gorm:"default:0"`
	UpdatedAt   int    `gorm:"default:0"`
}

const (
	AutoloadNo  = 0
	AutoloadYes = 1
)

func (Option) TableName() string { return "options" }
