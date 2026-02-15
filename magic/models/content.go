package models

type Content struct {
	Cid          int    `gorm:"primaryKey;autoIncrement"`
	Title        string `gorm:"size:150;not null"`
	Slug         string `gorm:"size:150;unique;not null"`
	Created      int    `gorm:"default:0"`
	Modified     int    `gorm:"default:0"`
	Text         string `gorm:"type:text"`
	Order        int    `gorm:"default:0"`
	AuthorID     int    `gorm:"not null;index"`
	Template     string `gorm:"size:32;default:''"`
	Type         string `gorm:"size:16;default:'post'"`
	Status       string `gorm:"size:16;default:'draft'"`
	Password     string `gorm:"size:32;default:''"`
	CommentsNum  int    `gorm:"default:0"`
	AllowComment bool   `gorm:"default:true"`
	AllowPing    bool   `gorm:"default:true"`
	AllowFeed    bool   `gorm:"default:true"`
	Parent       int    `gorm:"default:0"`
}

const (
	TypePost = "post"
	TypePage = "page"
)

const (
	StatusDraft   = "draft"
	StatusPublish = "publish"
	StatusPrivate = "private"
)

func (Content) TableName() string { return "content" }
