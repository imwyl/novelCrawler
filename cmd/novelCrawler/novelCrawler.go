package main

import (
	"strings"
	"flag"
	"log"

	"github.com/imwyl/novelCrawler/dao"

	"github.com/imwyl/novelCrawler/crwaler"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	var URL string
	flag.StringVar(&URL, "URL", "", "飘天小说地址")
	flag.Parse()
	if URL == "" {
		log.Fatalln("No URL")
		return
	} else if !strings.Contains(URL, "piaotian") {
		log.Fatalln("Incorrent URL")
	}
	db, err := dao.GetDB()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer db.Close()
	db.AutoMigrate(&dao.Chapter{}, &dao.Novel{})
	crwaler.Start(URL)
}
