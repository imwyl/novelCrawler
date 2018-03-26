package crwaler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

var r, _ = regexp.Compile("第[一二三四五六七八九零十百千万]+章")
var c = colly.NewCollector(
	colly.AllowedDomains("www.piaotian.com", "piaotian.com"),
	colly.CacheDir("/tmp/colly-cache"))

func initialize(URL string) {
	d := c.Clone()
	c.OnHTML("ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			if r.MatchString(el.ChildText("a")) {
				d.Visit(e.Request.AbsoluteURL(el.ChildAttr("a", "href")))
			}
		})
	})
	d.OnHTML("div#main", func(e *colly.HTMLElement) {
		fmt.Println(strings.Replace(e.Text, "\u00a0\u00a0\u00a0\u00a0", "<p>", -1))
	})
	d.OnResponse(func(res *colly.Response) {
		res.Body = []byte(strings.Replace(string(res.Body), "<br>", "<div id=\"main\">", 1))
	})
	c.Visit(URL)
	c.Wait()
	c.Limit(&colly.LimitRule{
		Parallelism: 5,
	})
}

// Start start the crawler
func Start() {
	initialize("https://www.piaotian.com/html/0/296/index.html")
}
