package kt

import (
	"testing"
)

func TestGetDSN(t *testing.T) {
	tests := []struct {
		name     string
		env      []string
		expected string
	}{
		{
			name:     "no variables",
			env:      []string{},
			expected: "",
		},
		{
			name:     "no matching variables",
			env:      []string{"foo=bar"},
			expected: "",
		},
		{
			name: "main key matches",
			env: []string{
				envPrefix + "=http://foo.com/",
			},
			expected: "http://foo.com/",
		},
		{
			name: "sub key matches",
			env: []string{
				envPrefix + "_16=http://foo.com/",
			},
			expected: "http://foo.com/",
		},
		{
			name: "main takes priority over sub key",
			env: []string{
				envPrefix + "_16=http://foo.com/",
				envPrefix + "=http://bar.com/",
			},
			expected: "http://bar.com/",
		},
		{
			name: "sub keys properly sorted",
			env: []string{
				envPrefix + "_16=http://foo.com/",
				envPrefix + "_21=http://bar.com/",
			},
			expected: "http://bar.com/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getDSN(test.env)
			if test.expected != result {
				t.Errorf("Unexpected dsn returned: %s", result)
			}
		})
	}
}
