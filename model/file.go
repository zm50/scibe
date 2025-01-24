package model

import (
	"gorm.io/gorm"
)

type DocFile struct {
	gorm.Model
	Uid uint
	Name string
	Extension string
	Md5 []byte `gorm:"type:bytea"`
}

func (DocFile) TableName() string {
	return "backend.doc_file"
}
