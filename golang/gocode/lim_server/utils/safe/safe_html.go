package safe

import (
	"regexp"
	"strings"
)

func SanitizeHTML(inputHTML string) string {
	re := regexp.MustCompile(`(?i)(</?script.*?>)|(</?iframe.*?>)|(onerror)|(</?meta.*?>)`)
	ok := re.MatchString(inputHTML)
	if ok {
		s := strings.ReplaceAll(inputHTML, "<", "&lt;")
		s = strings.ReplaceAll(s, ">", "&gt;")
		s = strings.ReplaceAll(s, "\"", "&quot;")
		return s
	}
	return inputHTML
}
