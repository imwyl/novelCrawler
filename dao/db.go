package dao

import (
	"github.com/jinzhu/gorm"
)


// Novel the novels 
type Novel struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null"`
	Chapters []Chapter
}

// Chapter chapters of novels
type Chapter struct {
	ID uint `gorm:"primary_key"`
	ChapterID string `gorm:"not null"`
	Name string 
	NovelID Novel
	Content string `grom:"type:text"`
}