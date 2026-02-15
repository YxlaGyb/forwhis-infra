package models

type Usermeta struct {
	UmetaId    int    `gorm:"primaryKey;autoIncrement"`
	UserId     int    `gorm:"not null;index"`
	MetaKey    string `gorm:"size:255;not null"`
	MetaValue  string `gorm:"type:text"`
	CreatedAt  int    `gorm:"default:0"`
	UpdatedAt  int    `gorm:"default:0"`
}

const (
	GenderSecret = 0
	GenderMale   = 1
	GenderFemale = 2
	GenderOther  = 3
)

const (
	MetaGender   = "gender"
	MetaBirthday = "birthday"
	MetaBio      = "bio"
)

func (Usermeta) TableName() string { return "usermeta" }
