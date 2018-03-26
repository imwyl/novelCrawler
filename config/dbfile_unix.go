// +build !windows

package novelcrawler

import (
	"os"
)

func getabspath() string {
	return getdbpath() + "/" + getdbname()
}

func getdbname() string {
	return ".novel.db"
}

func getdbpath() string {
	return os.Getenv("HOME")
}
