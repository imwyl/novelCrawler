package dao

import (
	"time"

	"github.com/imwyl/novelCrawler/config"
	"github.com/jinzhu/gorm"
)

// Novel the novels
type Novel struct {
	ID       string `grom:"primary_key; auto_increment:false"`
	Name     string `gorm:"type:varchar(100);not null"`
	Chapters []Chapter	`grom:"foreignkey:NovelID"`
	UpdateAt time.Time `gorm:"not null"`
}

// Chapter the chapters of volumn
type Chapter struct {
	ID      string `gorm:"primary_key; auto_increment:false"`
	Name    string
	NovelID string	`gorm:"not null"`
	Content string `gorm:"type:Text(100000)"`
	Orders  uint   `gorm:"not null"`
}

// GetDB returns a database connection
func GetDB() (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", novelcrawler.Getabspath())
	return
}

// Save a chapter entity
func (chapter *Chapter) Save() (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}
	if db.NewRecord(chapter) {
		db.Create(chapter)
	}
	db.Find(chapter)
	var novel Novel
	db.Model(chapter).Related(&novel)
	db.Save(&novel)
	return true, nil
}

// RecordExists check a record exists or not
func RecordExists(value interface{}) bool {
	db, err := GetDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var count int
	db.Find(value).Count(&count)
	return count > 0
}
