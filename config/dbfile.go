package novelcrawler

import (
	"log"
	"os"
)

// Getdbname returns the name of db file
func Getdbname() string {
	return getdbname()
}

// Getdbpath returns the path of db file 
func Getdbpath() string {
	path := getdbpath()
	_, err := os.Stat(path)
	if err != nil {
		result := os.MkdirAll(path, 777)
		if result != nil {
			log.Fatalf("Create dir %s failed", path)
			return ""
		}
		log.Println("Create dir:", path)
	}
	return path
}

// Getabspath return absolute path of db file
func Getabspath() string {
	return getabspath()
}

// GetTempDir returns temp directory on diffrent OS
func GetTempDir() string {
	return getTempDir()
}