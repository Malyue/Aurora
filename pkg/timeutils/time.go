package timeutils

import (
	"fmt"
	"github.com/golang-module/carbon"
	"strings"
)

// GetTodayDate return format as '2023-12-19'
func GetTodayDate() string {
	return carbon.Now().ToDateString()
}

// GetNowTime return format
func GetNowTime() string {
	timeStr := strings.Split(carbon.Now().String(), " ")
	return fmt.Sprintf("%s-%s", timeStr[0], timeStr[1])
}

// GetMonthDate return 2023-10
func GetMonthDate() string {
	year := carbon.Now().Year()
	month := carbon.Now().Month()
	return fmt.Sprintf("%d-%d", year, month)
}

// GetSeasonDate return 2023-Autumn
func GetSeasonDate() string {
	year := carbon.Now().Year()
	season := carbon.Now().Season()
	return fmt.Sprintf("%d-%s", year, season)
}
