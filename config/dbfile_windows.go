// +build windows 

package novelcrawler

import (
	"os"
)

func getdbname() string {
	return "novel.db"
}

func getdbpath() string {
	return os.Getenv("APPDATA") + "\\novelCrawler\\"
}

func getabspath() string {
	return Getdbpath() + getdbname()
}

func getTempDir() string {
	return os.Getenv("TEMP") + "\\colly-cache"
}