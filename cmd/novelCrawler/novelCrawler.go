package main

import (
	"log"
	"github.com/imwyl/novelCrawler/dao"
	"flag"
	"fmt"

	"github.com/imwyl/novelCrawler/crwaler"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var URL string
	flag.StringVar(&URL, "URL", "", "飘天小说地址")
	flag.Parse()
	if URL == "" {
		fmt.Println("No URL")
	} else {
		crwaler.Start(URL)
	}
	db, err := dao.GetDB()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()
	db.AutoMigrate(&dao.Chapter{}, &dao.Novel{})
}
