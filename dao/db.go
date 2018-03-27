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
	Content string `grom:"type:text"`
	Orders  uint   `gorm:"not null"`
}

// ExistQuery an interface of dao
type ExistQuery interface {
	// Exists returns whether a entity exists or not
	Exists() (bool, error)
}

// GetDB returns a database connection
func GetDB() (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", novelcrawler.Getabspath())
	return
}

// Exists returns whether a chapter exists or not
func (chapter *Chapter) Exists() (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}
	defer db.Close()
	return !db.NewRecord(chapter), nil
}

// Save a chapter entity
func (chapter *Chapter) Save() (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}
	if exits, err := chapter.Exists(); !exits && err != nil {
		db.Create(chapter)
	}
	db.Find(chapter)
	var novel Novel
	db.Model(chapter).Related(&novel)
	db.Save(&novel)
	return true, nil
}

// Exists returns whether a novel exists or not
func (novel *Novel) Exists() (bool, error) {
	db, err := GetDB()
	if err != nil {
		return false, err
	}
	defer db.Close()
	return !db.NewRecord(novel), nil
}
