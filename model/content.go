package model

import "gorm.io/gorm"

type Content struct {
	gorm.Model
	Fid uint
	Extraction string
	FileMindmap string
	FileSummary string
}

func (Content) TableName() string {
	return "backend.content"
}
