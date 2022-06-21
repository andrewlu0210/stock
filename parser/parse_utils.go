package parser

import (
	"fmt"
	"strconv"
	"time"
)

//日期格式錯誤或日期是未來日期(after today)會返回error
func ParseDate(dateStr string) (time.Time, error) {
	layout := "20060102"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return date, err
	}
	now := time.Now()
	if date.After(now) {
		return date, fmt.Errorf(fmt.Sprintf("The input date is future date:%s", dateStr))
	}
	return date, nil
}

func StringToInt(word string) int {
	val, err := strconv.ParseInt(word, 10, 64)
	if err != nil {
		return 0
	}
	return int(val)
}

func StringToLong(word string) int64 {
	val, err := strconv.ParseInt(word, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

func StringToFloat(word string) float64 {
	val, err := strconv.ParseFloat(word, 64)
	if err != nil {
		return 0
	}
	return val
}
