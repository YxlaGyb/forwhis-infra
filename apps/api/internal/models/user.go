package models

type User struct {
	Uid       int    `gorm:"primaryKey"`
	Name      string `gorm:"size:32;uniqueIndex"`
	Mail      string `gorm:"size:150;uniqueIndex"`
	Password  string `gorm:"size:255"`
	Gender    int    `gorm:"default:0"`
	Birthday  string `gorm:"size:10;default:''"`
	Bio       string `gorm:"size:200;default:''"`
	CreatedAt int    `gorm:"default:0"`
	LastLogin int    `gorm:"default:0"`
	LastIp    string `gorm:"size:150;default:''"`
	IsActive  int    `gorm:"default:0"`
}

const (
	GenderSecret = 0
	GenderMale   = 1
	GenderFemale = 2
	GenderOther  = 3
)

const (
	ProfileGender   = "gender"
	ProfileBirthday = "birthday"
	ProfileBio      = "bio"
)

func (User) TableName() string { return "users" }
