package novelcrawler

import (
	"fmt"
)

// Chstonum returns the uint number of chinese number, 
// if unexpected number presents, or larger than 99999999, 
// an error will be returned
func Chstonum(strnum string) (result int, err error) {
	tmp := -1
	for _, ch := range strnum {
		pos, err := positionOf(ch)	
		if err != nil {
			return pos, err
		}
		if pos == 0 {	// 零
			tmp = -1
		} else if pos > 0 && pos < 10{ // 一 ~ 九
			tmp = pos
		} else if pos >= 10 && pos < 10000{ // 十， 百， 千
			if tmp <= 0 && result == 0{
				tmp = 1
			}
			result += tmp * pos
			tmp = -1
		} else { // 万
			// if the number before 10000(万) is zero, it's a illegal number
			if result == 0 && tmp <= 0{
				return -1, fmt.Errorf("Illegal number %v", strnum)
			} else if tmp > 0 {
				result += tmp
			}
			result *= pos
		}
	}
	if tmp != -1 {
		result += tmp
	}
	return
}

// positionOf return the position of a rune, 
// if not in the slice, return an error
func positionOf(ch rune) (int, error) {
	if ch == '两' {
		ch = '二'
	}
	chsnums := []rune{'零', '一', '二', '三', '四', '五', '六', '七', '八', '九'}
	for pos, value := range chsnums {
		if ch == value {
			return pos, nil
		}
	}
	return bits(ch)
}

func bits(ch rune) (result int, err error) {
	switch ch {
	case '十':
		result = 10
	case '百':
		result = 100
	case '千':
		result = 1000
	case '万':
		result = 10000
	default:
		return -1, fmt.Errorf("Unexpcted chinese number: %v", ch)
	}
	return
}