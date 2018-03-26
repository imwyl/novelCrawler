package novelcrawler

import (
	"github.com/imwyl/novelCrawler/pkg/novelCrawler"
	"testing"
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
