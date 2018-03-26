package novelcrawler

import (
	nc "github.com/imwyl/novelCrawler/config"
	"testing"
)

func TestGetdbname(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test1", ".novel.db"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nc.Getdbname(); got != tt.want {
				t.Errorf("Getdbname() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetdbpath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test2", "/home/wyl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nc.Getdbpath(); got != tt.want {
				t.Errorf("Getdbpath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetabspath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test3", "/home/wyl/.novel.db"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := nc.Getabspath(); got != tt.want {
				t.Errorf("Getabspath() = %v, want %v", got, tt.want)
			}
		})
	}
}
