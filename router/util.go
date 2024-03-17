package router

import (
	"fmt"
	"regexp"
)

func match(re *regexp.Regexp, target string) (bool, map[string]string) {
	matches := re.FindStringSubmatch(target)

	if len(matches) == 0 {
		return false, nil
	}

	groups := make(map[string]string)
	names := re.SubexpNames()

	for i, v := range matches[1:] {
		i++
		if names[i] != "" {
			groups[names[i]] = v
		} else {
			groups[fmt.Sprint(i)] = v
		}
	}

	return true, groups
}
