package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	OriginalFilename string
	Md5 string
}
