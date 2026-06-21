package power

import "regexp"

// 预编译正则表达式，避免在循环中重复编译
var (
	reFloatPower = regexp.MustCompile(`([0-9]+\.[0-9]+)\s*W`)
	reIntPower   = regexp.MustCompile(`([0-9]+)\s*W`)
)

// regexpMustCompile 包装 regexp.MustCompile，便于测试
func regexpMustCompile(pattern string) *regexp.Regexp {
	switch pattern {
	case `([0-9]+\.[0-9]+)\s*W`:
		return reFloatPower
	case `([0-9]+)\s*W`:
		return reIntPower
	}
	return regexp.MustCompile(pattern)
}
