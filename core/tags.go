package core

import (
	"slices"
	"strings"
)

type TagOptions []string

func ParseTag(tag string) (string, TagOptions) {
	ret := strings.Split(tag, ",")
	return ret[0], ret[1:]
}

func (o TagOptions) Contains(option string) bool {
	return slices.Contains(o, option)
}
