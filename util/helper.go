package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func EnforceInt(value interface{}) int {

	switch v := value.(type) {

	case float64:
		return int(v)

	case int64:
		return int(v)

	case int32:
		return int(v)

	case string:

		k, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return 0
		}

		return int(k)

	case int:
		return v

	}

	return 0

}

func GetIntegerFromString(s string) int {

	val, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return val

}

func NewValidationError(message string) error {
	return fmt.Errorf(message)
}

func FormatError(err error) string {

	s := err.Error()
	s = strings.ToUpper(string(s[0])) + s[1:] + "."

	return s

}

func InList(s string, list []string) bool {

	for _, item := range list {
		if item == s {
			return true
		}
	}

	return false

}

func RemoveEmptyStrings(list []string) []string {

	var newList []string

	for _, item := range list {

		item = strings.TrimSpace(item)

		if item != "" {
			newList = append(newList, item)
		}

	}

	return newList

}

func GetStringFromInteger(i int) string {

	return strconv.Itoa(i)

}

func FormatWithOrdinalSuffix(t time.Time) string {

	day := t.Day()

	var suffix string

	if day%10 == 1 && day != 11 {
		suffix = "st"
	} else if day%10 == 2 && day != 12 {
		suffix = "nd"
	} else if day%10 == 3 && day != 13 {
		suffix = "rd"
	} else {
		suffix = "th"
	}

	return fmt.Sprintf("%d%s %s %d", day, suffix, t.Format("January"), t.Year())

}

func FormatDate(dateStr string) string {

	t, err := time.Parse("2006-01-02", dateStr)

	if err != nil {
		return ""
	}

	return FormatWithOrdinalSuffix(t)

}

func GetDateString(t time.Time) string {

	return t.Format("2006-01-02")

}

func GetHourString(t time.Time) string {

	return t.Format("2006-01-02 15:04")

}

func GetClampedPageNumber(page int, totalPages int) int {

	if page < 1 {
		return 1
	}

	if page > totalPages {
		return totalPages
	}

	return page

}

func GetTotalPages(totalItems int, itemsPerPage int) int {

	if totalItems == 0 {
		return 1
	}

	return (totalItems + itemsPerPage - 1) / itemsPerPage

}
