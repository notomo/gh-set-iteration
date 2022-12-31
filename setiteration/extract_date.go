package setiteration

import (
	"fmt"
	"regexp"
)

var dateRegex = regexp.MustCompile(`\d\d\d\d-\d\d-\d\d`)

func ExtractDate(title string) (string, error) {
	date := dateRegex.FindString(title)
	if date == "" {
		return "", fmt.Errorf("title should include yyyy-mm-dd: title='%s'", title)
	}
	return date, nil
}
