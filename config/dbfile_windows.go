// +build windows 

package novelcrawler

import (
	"os"
)

func getdbname() string {
	return "novel.db"
}

func getdbpath() string {
	return os.Getenv("APPDATA")
}

func getabspath() string {
	return getdbpath() + "\\" + getdbname()
}