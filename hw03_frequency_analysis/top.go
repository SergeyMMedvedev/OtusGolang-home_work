package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var wordRegex = regexp.MustCompile(
	`([а-яА-Я\w]+(?:[-\.,]+[-а-яА-Я\w]+)?)|-{2,}`,
)

func Top10(str string) (res []string) {
	fields := wordRegex.FindAllString(strings.ToLower(str), -1)
	freq := make(map[string]int)
	arr := make([][]string, len(fields))
	for _, word := range fields {
		freq[word]++
	}
	for k, v := range freq {
		arr[v] = append(arr[v], k)
	}
	for i := len(arr) - 1; i >= 0; i-- {
		if len(arr[i]) > 0 {
			sort.Strings(arr[i])
			res = append(res, arr[i]...)
		}
	}
	if len(res) > 10 {
		return res[:10]
	}
	return res
}
