package dao

import (
	"github.com/jinzhu/gorm"
)


// Novel the novels 
type Novel struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);not null"`
	URL string `gorm:"not null"`
	Volumns []Volumn
}

// Volumn the volumns of a novel
type Volumn struct {
	ID uint `gorm:"primary_key"`
	Name string 
	order uint `gorm:"not null"`
	NovelID Novel `gorm:"not null"`
	Chapters []Chapter
}

// Chapter the chapters of volumn
type Chapter struct {
	ID uint `gorm:"primary_key"`
	ChapterID string `gorm:"not null"`
	Name string 
	URL string `gorm:"not null"`
	VolumnID Volumn `gorm:"not null"`
	Content string `grom:"type:text"`
	order uint `gorm:"not null"`
}