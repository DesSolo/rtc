package generator

import (
	"strings"
	"text/template"
	"unicode"
)

// ParseWithFunctions ...
func ParseWithFunctions(text string) (*template.Template, error) {
	return template.New("").Funcs(advancedFunctions).Parse(text) // nolint:wrapcheck
}

var advancedFunctions = template.FuncMap{
	"toPascal": func(s string) string {
		var result strings.Builder

		for _, word := range strings.Split(s, "_") {
			if word == "" {
				continue
			}

			runes := []rune(word)
			first := unicode.ToUpper(runes[0])

			rest := string(runes[1:])
			if len(rest) > 0 {
				rest = strings.ToLower(rest)
			}

			result.WriteRune(first)
			result.WriteString(rest)
		}

		return result.String()
	},
}
