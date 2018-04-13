package crwaler

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/imwyl/novelCrawler/config"
	"github.com/imwyl/novelCrawler/pkg/novelCrawler"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
	"github.com/imwyl/novelCrawler/dao"
)

var chapterRegex, _ = regexp.Compile("(第[一二三四五六七八九零十百千万]+[章|节]\\s+.*)")
var novelNameRegex, _ = regexp.Compile("^(.*)最新章节$")
var novelURIRegex, _ = regexp.Compile(`\d+/\d+`)
var chapterURLRegex, _ = regexp.Compile("(\\d+)\\.html")
var firstChapterRegex, _ = regexp.Compile(`<a href="(\d+).html">第一[章|节]`)
var tempDir = config.GetTempDir()

func getChapterNumber(chapterURL string) int {
	if chapterURL == "" {
		return -1
	}
	chapterNumberArr := chapterURLRegex.FindStringSubmatch(chapterURL)
	if len(chapterNumberArr) <= 1 {
		log.Fatalln(chapterURL)
		return -1
	}
	result, err := strconv.Atoi(chapterNumberArr[1])
	if err != nil {
		return -1
	}
	return result
}

func getChapters(URL string) {
	var c = colly.NewCollector(
		colly.AllowedDomains("piaotian.net", "www.piaotian.net", "www.piaotian.com", "piaotian.com" ),
		colly.CacheDir(tempDir))
	c.Limit(&colly.LimitRule{
		Parallelism: 5,
	})
	enableProxy(c)
	d := c.Clone()
	enableProxy(d)
	d.Limit(&colly.LimitRule{
		Parallelism: 5,
	})
	db, err := dao.GetDB()
	if err != nil {
		log.Fatalln(err)
	}
	var chapter dao.Chapter
	novel := dao.Novel{ID: getNovelURI(URL)}
	db.First(&novel)
	db.Where("novel_id = ?", novel.ID).Order("id desc").First(&chapter)
	chapterMap := make(map[int]*dao.Chapter)
	// get chapters of novel
	c.OnHTML("ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			if chapterRegex.MatchString(el.ChildText("a")) {
				chapterURL := el.ChildAttr("a", "href")
				if thisChapter := getChapterNumber(chapterURL); thisChapter-novel.First+1 > int(chapter.ID) {
					if chapterMap[thisChapter] == nil {
						chapterMap[thisChapter] = &dao.Chapter{
							ID:      uint(thisChapter),
							NovelID: novel.ID,
							Name:    chapterRegex.FindStringSubmatch(el.ChildText("a"))[1],
						}
						log.Printf("Add %d %v to map", thisChapter, chapterMap[thisChapter])
					}
					d.Visit(e.Request.AbsoluteURL(el.ChildAttr("a", "href")))
				}
			}
		})
	})

	// get content of novel
	d.OnHTML("div#main", func(e *colly.HTMLElement) {
		requestURL := chapterURLRegex.FindStringSubmatch(e.Request.URL.String())[1]
		chapterURL := getChapterNumber(requestURL)
		thisChapter := chapterMap[chapterURL]
		log.Println("div#main", thisChapter)
		thisChapter.Content = strings.Replace(e.Text, "\u00a0\u00a0\u00a0\u00a0", "<p>", -1)
		thisChapter.Save()
		//log.Println()
	})

	// div#main is add by javascript, so add it here
	d.OnResponse(func(res *colly.Response) {
		res.Body = []byte(strings.Replace(string(res.Body), "<br>", "<div id=\"main\">", 1))
		log.Println("Response code: ", string(res.Body))
	})

	c.Visit(URL)
	c.Wait()

}

// novelExits return the novel exits in database or not
func novelExists(URL string) bool {
	novelURI := getNovelURI(URL)
	log.Println(novelURI)
	return dao.NovelExists(novelURI)
}

// createNovel create novel record in database
func createNovel(URL string, boolChan chan novelcrawler.ErrorCode) {
	log.Println("URL", URL)
	db, getDbErr := dao.GetDB()
	if getDbErr != nil {
		log.Fatalln("Can not open databse:\n", getDbErr)
		panic(getDbErr)
	}
	defer db.Close()
	var c = colly.NewCollector(
		colly.AllowedDomains("www.piaotian.com", "piaotian.com", "piaotian.net", "www.piaotian.net"),
		colly.CacheDir(tempDir))
	enableProxy(c)
	var remoteNovelExist bool
	var firstChapter int
	c.OnResponse(func(res *colly.Response) {
		remoteNovelExist = res.StatusCode == 200
		if !remoteNovelExist {
			boolChan <- novelcrawler.RemoteNotExist
		} else {
			firstChapterSlice := firstChapterRegex.FindStringSubmatch(string(res.Body))
			if len(firstChapterSlice) <= 1 {
				boolChan <- novelcrawler.InternalError
			} else {
				firstChapter, _ = strconv.Atoi(firstChapterSlice[1])
			}
		}
	})
	c.OnHTML("h1", func(e *colly.HTMLElement) {
		h1title := e.Text
		// get novel's name and save to database
		if novelName := novelNameRegex.FindStringSubmatch(h1title); len(novelName) > 1 && remoteNovelExist {
			log.Println(novelName)
			novel := dao.Novel{
				ID:       getNovelURI(URL),
				Name:     novelName[1],
				First:    firstChapter,
				UpdateAt: time.Now(),
			}
			log.Println("create novel ", novel.Name)
			db.Create(&novel)
			boolChan <- novelcrawler.Success
		} else {
			boolChan <- novelcrawler.NovelName
			log.Println("novelExists ", remoteNovelExist)
		}
	})
	c.Visit(URL)
	c.Wait()
}

// Start the crawler
func Start(URL string) novelcrawler.ErrorCode {
	modifyURL(&URL)
	if !novelExists(URL) {
		errorChan := make(chan novelcrawler.ErrorCode)
		go createNovel(URL, errorChan)
		if createResult := <-errorChan; createResult != novelcrawler.Success {
			return createResult
		}
	}
	getChapters(URL)
	return novelcrawler.Success
}

func modifyURL(URL *string) {
	if !strings.Contains(*URL, "https://www.") {
		*URL = "https://www." + *URL
	}
	if strings.HasSuffix(*URL, "/") {
		*URL = string((*URL)[:len(*URL)-1])
	} else {
		*URL = strings.Replace(*URL, "/index.html", "", 1)
	}
}

func enableProxy(c *colly.Collector) {
	p, proxyErr := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:9050")
	if proxyErr != nil {
		panic(proxyErr)
	}
	c.SetProxyFunc(p);
}
func getNovelURI(URL string) string {
	return novelURIRegex.FindString(URL)
}
