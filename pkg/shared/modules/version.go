package modules

import (
	"fmt"
	"strings"
)

type Version [3]int

func NewVersion(major int, minor int, sub int) Version {
	return Version{major, minor, sub}
}

func (v Version) string() string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v)), "."), "[]")
}
