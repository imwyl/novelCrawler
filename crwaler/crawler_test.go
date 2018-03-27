package crwaler

import "testing"

func Test_chapterOrder(t *testing.T) {
	type args struct {
		chapterURL string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"1", args{"1234221.html"}, 1234221},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getChapterOrder(tt.args.chapterURL); got != tt.want {
				t.Errorf("chapterOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
