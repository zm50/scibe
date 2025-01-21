package model

import "gorm.io/gorm"

type Content struct {
	gorm.Model
	FilePath string
	FileExtraction string
	FileMindmap string
	FileSummary string
}
