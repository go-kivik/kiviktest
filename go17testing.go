// +build go1.7,!go1.8

package kiviktest

import (
	"os"
	"regexp"
	"testing"
)

func mainStart(tests []testing.InternalTest) {
	m := testing.MainStart(regexp.MatchString, tests, nil, nil)
	os.Exit(m.Run())
}
