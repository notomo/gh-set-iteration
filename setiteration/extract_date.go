package setiteration

import (
	"fmt"
	"regexp"
	"time"
)

var dateRegex = regexp.MustCompile(`\d\d\d\d-\d\d-\d\d`)

func ExtractDate(title string) (string, error) {
	date := dateRegex.FindString(title)
	if date == "" {
		return "", fmt.Errorf("title should include yyyy-mm-dd: title='%s'", title)
	}
	return date, nil
}

const dateFormat = "2006-01-02"

func ShiftDate(date string, offsetDays int) (string, error) {
	at, err := time.Parse(dateFormat, date)
	if err != nil {
		return "", err
	}
	startDate := at.AddDate(0, 0, offsetDays)
	return startDate.Format(dateFormat), nil
}
