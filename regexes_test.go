package porto

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegexpList(t *testing.T) {
	tests := []struct {
		name     string
		regexps  string
		expected []*regexp.Regexp
		err      string
	}{
		{
			name:     "single regex",
			regexps:  "^.*pb.go$",
			expected: []*regexp.Regexp{regexp.MustCompile("^.*pb.go$")},
		},
		{
			name:    "failing regex",
			regexps: "^.*pb.go$,*$$",
			err:     "failed to compile regex \"*$$\": error parsing regexp: missing argument to repetition operator: `*`",
		},
		{
			name:     "multiple regexes",
			regexps:  "^.*pb.go$,^tools.go$",
			expected: []*regexp.Regexp{regexp.MustCompile("^.*pb.go$"), regexp.MustCompile("^tools.go$")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := GetRegexpList(tt.regexps)
			if err != nil || tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
