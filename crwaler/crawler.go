package crwaler

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/imwyl/novelCrawler/dao"
)

var chapterRegex, _ = regexp.Compile("第([一二三四五六七八九零十百千万]+)[章|节]\\s(.*)")
var titleRegex, _ = regexp.Compile("^(.*)最新章节$")
var chapterURLRegex, err = regexp.Compile("(\\d+)")

func getChapterOrder(chapterURL string) int {
	if chapterURL == "" {
		return -1
	}
	chapterNumber := chapterURLRegex.FindString(chapterURL)
	result, err := strconv.Atoi(chapterNumber)
	if err != nil {
		return -1
	}
	return result
}

func initialize(URL string) {
	var c = colly.NewCollector(
		colly.AllowedDomains("www.piaotian.com", "piaotian.com"),
		colly.CacheDir("/tmp/colly-cache"))
	d := c.Clone()

	db, err := dao.GetDB()
	if err != nil {
		log.Fatalln(err)
	}
	var chapter dao.Chapter
	db.Where("novel_id = ?", URL).Order("orders desc").First(&chapter)
	chapterNumber := getChapterOrder(chapter.ID)
	fmt.Println(chapterNumber)
	// get chapters of novel
	// c.OnHTML("ul", func(e *colly.HTMLElement) {
	// 	fmt.Println(chapter)
	// 	e.ForEach("li", func(_ int, el *colly.HTMLElement) {
	// 		if chapterRegex.MatchString(el.ChildText("a")) {
	// 			chapterURL := el.ChildAttr("a", "href")
	// 			if getChapterOrder(chapterURL) > chapterNumber {
	// 				log.Printf("Visiting %s\n", chapterURL)
	// 				// d.Visit(e.Request.AbsoluteURL(el.ChildAttr("a", "href")))
	// 			}
	// 		}
	// 	})
	// })

	// get content of novel
	d.OnHTML("div#main", func(e *colly.HTMLElement) {
		//fmt.Println(strings.Replace(e.Text, "\u00a0\u00a0\u00a0\u00a0", "<p>", -1))
	})

	// get volume and chapter of novel
	d.OnHTML("h1", func(e *colly.HTMLElement) {
		fmt.Println(chapterRegex.FindStringSubmatch(e.Text))
	})

	// div#main is add by javascript, so add it here
	d.OnResponse(func(res *colly.Response) {
		res.Body = []byte(strings.Replace(string(res.Body), "<br>", "<div id=\"main\">", 1))
	})

	c.Visit(URL)
	c.Wait()

	c.Limit(&colly.LimitRule{
		Parallelism: 5,
	})
}

// novelExits return the novel exits in database or not
func novelExits(URL string, boolChan chan bool) {
	db, err := dao.GetDB()
	if err != nil {
		log.Fatalln("Can not open databse:\n", err)
		panic(err)
	}
	URL = strings.Replace(URL, "/index.html", "", 1)
	fmt.Println(URL)
	novel := dao.Novel{ID: URL}
	defer db.Close()
	novelExist, err := novel.Exists()
	fmt.Println("Database exists", novelExist)
	if err != nil {
		log.Fatalln(err)
		boolChan <- false
		panic(err)
	}
	if novelExist {
		boolChan <- true
		return
	}
	var c = colly.NewCollector(
		colly.AllowedDomains("www.piaotian.com", "piaotian.com"),
		colly.CacheDir("/tmp/colly-cache"))
	c.OnResponse(func(res *colly.Response) {
		log.Println("Response recived")
		novelExist = res.StatusCode == 200
	})
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		h1title := e.Text
		// get novel name and save to database
		if novelName := titleRegex.FindStringSubmatch(h1title); len(novelName) != 0 && novelExist {
			fmt.Println(novelName)
			novel.Name = novelName[1]
			novel.ID = URL
			novel.UpdateAt = time.Now()
			log.Println("create novel ", novel.Name)
			db.Create(&novel)
			boolChan <- novelExist
		} else {
			if err != nil {
				log.Fatalln(err)
				boolChan <- false
			}
			log.Println("novelExists ", novelExist)
		}
	})
	c.Visit(URL)
	c.Wait()
}

// Start the crawler
func Start(URL string) {
	boolChan := make(chan bool)
	go novelExits(URL, boolChan)
	if exists := <-boolChan; exists {
		initialize(URL)
	}
}
