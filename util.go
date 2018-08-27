package gosaic

import (
	"strings"
)

func SplitPath(path string) []string {
	trimmed := strings.TrimFunc(path, func(r rune) bool {
		return r == '/'
	})

	return strings.Split(trimmed, "/")
}
