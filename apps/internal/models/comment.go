package models

type Comment struct {
	Coid     int    `gorm:"primaryKey;autoIncrement"`
	Cid      int    `gorm:"not null;index"`
	Created  int    `gorm:"default:0"`
	Author   string `gorm:"size:150;default:''"`
	AuthorID int    `gorm:"index"`
	OwnerID  int    `gorm:"index"`
	Mail     string `gorm:"size:150;default:''"`
	Url      string `gorm:"size:255;default:''"`
	Ip       string `gorm:"size:64;default:''"`
	Agent    string `gorm:"size:511;default:''"`
	Text     string `gorm:"type:text"`
	Type     string `gorm:"size:16;default:'comment'"`
	Status   string `gorm:"size:16;default:'pending'"`
	Parent   int    `gorm:"default:0"`
}

const (
	CommentTypeComment   = "comment"
	CommentTypePing      = "ping"
	CommentTypeTrackback = "trackback"
)

const (
	CommentStatusPending  = "pending"
	CommentStatusApproved = "approved"
	CommentStatusSpam     = "spam"
)

func (Comment) TableName() string { return "comment" }
