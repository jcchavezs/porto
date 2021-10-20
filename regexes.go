package porto

import (
	"fmt"
	"regexp"
	"strings"
)

func GetRegexpList(regexps string) ([]*regexp.Regexp, error) {
	var regexes []*regexp.Regexp
	if len(regexps) > 0 {
		for _, sfrp := range strings.Split(regexps, ",") {
			sfr, err := regexp.Compile(sfrp)
			if err != nil {
				return nil, fmt.Errorf("failed to compile regex %q: %w", sfrp, err)
			}
			regexes = append(regexes, sfr)
		}
	}

	return regexes, nil
}
