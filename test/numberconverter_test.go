package novelcrawler

import (
	"testing"

	"github.com/imwyl/novelCrawler/pkg/novelCrawler"
)

func TestChstonum(t *testing.T) {
	type args struct {
		strnum string
	}
	tests := []struct {
		name       string
		args       args
		wantResult int
		wantErr    bool
	}{
		{"0", args{"零"},
			0, false},
		{"1", args{"一"},
			1, false},
		{"2", args{"十二"},
			12, false},
		{"3", args{"一百零三"},
			103, false},
		{"3", args{"一零三"},
			103, false},
		{"4", args{"一百四十五"},
			145, false},
		{"5-1", args{"六千零七十八"},
			6078, false},
		{"5-3", args{"六千零一"},
			6001, false},
		{"5-2", args{"六零七八"},
			6078, false},
		{"6", args{"九千一百零二"},
			9102, false},
		{"7", args{"三千四百五十六"},
			3456, false},
		{"8", args{"七万零八"},
			70008, false},
		{"9", args{"九万零一百零二"},
			90102, false},
		{"10", args{"三万零四百五十六"},
			30456, false},
		{"11", args{"七万八千零九"},
			78009, false},
		{"12", args{"十万两千三百零四"},
			102304, false},
		{"13", args{"五十万六千七百八十九"},
			506789, false},
		{"13-1", args{"五零六七八九"},
			506789, false},
		{"14", args{"一百二十三万四千五百六十七"},
			1234567, false},
		{"err", args{"第一章"},
			-1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := novelcrawler.Chstonum(tt.args.strnum)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chstonum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("Chstonum() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestNumberToChs(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{1200}, "一千二百"},
		{"2", args{10}, "一十"},
		{"3", args{14}, "一十四"},
		{"4", args{20}, "二十"},
		{"5", args{100}, "一百"},
		{"6", args{104}, "一百零四"},
		{"7", args{112}, "一百一十二"},
		{"8", args{130}, "一百三十"},
		{"9", args{1000}, "一千"},
		{"10", args{1002}, "一千零二"},
		{"11", args{1030}, "一千零三十"},
		{"12", args{4056}, "四千零五十六"},
		{"13", args{7890}, "七千八百九十"},
		{"14", args{1200}, "一千二百"},
		{"15", args{102345}, "一十万零二千三百四十五"},
		{"16", args{40000}, "四万"},
		{"17", args{500000}, "五十万"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := novelcrawler.NumberToChs(tt.args.number); got != tt.want {
				t.Errorf("NumberToChs() = %v, want %v", got, tt.want)
			}
		})
	}
}
