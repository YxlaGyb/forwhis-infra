package models

type ContentField struct {
	Cid        int     `gorm:"primaryKey"`
	Name       string  `gorm:"size:150;primaryKey"`
	Type       string  `gorm:"size:8;default:'str'"`
	StrValue   string  `gorm:"type:text"`
	IntValue   int     `gorm:"default:0"`
	FloatValue float64 `gorm:"default:0"`
}

const (
	FieldTypeString = "str"
	FieldTypeInt    = "int"
	FieldTypeFloat  = "float"
)

func (ContentField) TableName() string { return "contentField" }
