package models

type Content struct {
	Cid        int    `gorm:"primaryKey;autoIncrement"`
	Title      string `gorm:"size:150;not null"`
	Created    int    `gorm:"default:0;index:idx_contents_feed,priority:3"`
	Modified   int    `gorm:"default:0"`
	Text       string `gorm:"type:text"`
	AuthorID   int    `gorm:"not null;index"`
	CategoryID int    `gorm:"default:0;index"`
	Type       string `gorm:"size:16;default:'post';index:idx_contents_feed,priority:2"`
	Status     string `gorm:"size:16;default:'draft';index:idx_contents_feed,priority:1"`
	Password   string `gorm:"size:32;default:''"`
	Views      int    `gorm:"default:0"`
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

func (Content) TableName() string { return "contents" }
