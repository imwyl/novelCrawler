package novelcrawler

import (
	"fmt"
	"log"
	"strings"
)

// Chstonum returns the int number of chinese number,
// if unexpected number presents, or larger than 99999999,
// an error will be returned
func Chstonum(strnum string) (result int, err error) {
	var method func(int, *int)
	var containsBits bool
	if strings.Contains(strnum, "百") ||
		strings.Contains(strnum, "千") ||
		strings.Contains(strnum, "万") ||
		strings.Contains(strnum, "十") {
		method = withbits()
		containsBits = true
	} else {
		method = withoutbits()
		containsBits = false
	}
	pos := 0
	for _, ch := range strnum {
		pos, err = positionOf(ch)
		if err != nil {
			return pos, err
		}
		method(pos, &result)
		if result == -1 {
			log.Fatalf("Error: unexpcted number: %v(%v)", strnum, ch)
			return -1, fmt.Errorf("Error: unexpcted number: %v(%v)", strnum, ch)
		}
	}
	if containsBits && pos != 0 {
		result += pos
	}
	return
}

func withbits() func(int, *int) {
	tmp := -1
	return func(pos int, result *int) {
		if pos == 0 {
			tmp = -1
		} else if pos > 0 && pos < 10 {
			tmp = pos
		} else if pos >= 10 && pos < 10000 {
			if tmp < 0 && *result <= 0 {
				tmp = 1
			}
			*result += (tmp * pos)
			tmp = -1
		} else {
			if *result == 0 && tmp <= 0 {
				*result = -1
			} else if tmp > 0 {
				*result += tmp
			}
			*result *= pos
			tmp = -1
		}
	}
}

func withoutbits() func(int, *int) {
	return func(pos int, result *int) {
		if pos == 0 {
			*result *= 10
		} else if pos > 0 && pos < 10 {
			*result *= 10
			*result += pos
		} else {
			*result = -1
		}
	}
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
	return bitToNumber(ch)
}

func bitToNumber(ch rune) (result int, err error) {
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

// NumberToChs converts int number to chinese number
func NumberToChs(number int) string {
	chs := []string{"", "十", "百", "千", "万"}
	chsnums := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	var result string
	current := -1
	last := -1
	for index := 0; number > 0; index = (index + 1) % len(chs) {
		if current != -1 {
			last = current
			if index == 0 {
				index = 1
			}
		}
		current = number % 10
		number /= 10
		if current == 0 && (last == -1 || last == 0) {
			if index == 4 {
				result += chs[index]
			}
			continue
		} else if current == 0 && index != 4 {
			result = chsnums[current] + result
		} else if current == 0 && index == 4 {
			result = chs[index] + chsnums[current] + result
		} else {
			result = chsnums[current] + chs[index] + result
		}
	}
	return result
}
