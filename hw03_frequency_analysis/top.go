package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var re = regexp.MustCompile(`[\s";.,!']+`)

func Top10(str string) (result []string) {
	if str == "" {
		return
	}

	words := re.Split(str, -1)
	if len(words) == 0 {
		return
	}

	m := make(map[string]int)
	for _, word := range words {
		w := strings.ToLower(word)

		if w == "-" {
			continue
		}

		m[w]++
	}

	keys := make([]string, 0, len(words))
	for word := range m {
		keys = append(keys, word)
	}

	sort.Slice(keys, func(i, j int) bool {
		if m[keys[i]] == m[keys[j]] {
			return keys[i] < keys[j]
		}

		return m[keys[i]] > m[keys[j]]
	})

	take := len(keys)
	if take > 10 {
		take = 10
	}

	result = make([]string, 0, take)
	for _, word := range keys {
		if len(result) >= take {
			break
		}
		result = append(result, word)
	}

	return
}
