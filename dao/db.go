package dao

import (
	"time"

	"github.com/imwyl/novelCrawler/config"
	"github.com/jinzhu/gorm"
)

// Novel the novels
type Novel struct {
	ID			string		`grom:"primary_key; auto_increment:false"`
	Name		string		`gorm:"type:varchar(100);not null"`
	Chapters	[]Chapter	`gorm:"foreignkey:NovelID"`
	First		int			`gorm:"not null"`
	UpdateAt	time.Time	`gorm:"not null"`
}

// Chapter the chapters of volumn
type Chapter struct {
	ID      uint	`gorm:"primary_key;auto_increment:false"`
	Name    string
	NovelID string `gorm:"not null"`
	Content string `gorm:"type:Text(100000)"`
}

// GetDB returns a database connection
func GetDB() (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", config.Getabspath())
	return
}

// Save a chapter entity
func (chapter *Chapter) Save() (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}
	var count int
	if db.Model(chapter).Where("id = ?", chapter.ID).Count(&count); count == 0 {
		db.Create(chapter)
	} else {
		db.Save(chapter)
	}
	var novel Novel
	db.Model(chapter).Related(&novel)
	db.Model(&novel).Update("update_at", time.Now())
	return true, nil
}

// NovelExists check a record exists or not
func NovelExists(id string) bool {
	db, err := GetDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var count int
	db.Model(&Novel{}).Where("id = ?", id).Count(&count)
	return count > 0
}
