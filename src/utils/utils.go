package utils

import (
	"fmt"
)

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func BoolToInt(b bool) int {
	switch b {
	case true:  return 1
	default:    return 0
	}
}

func RepeatString(s string, n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += s
	}
	return res
}

func Fmap[T any, F any](a []T, f func(T) F) []F {
	fs := make([]F, len(a))
	for i, v := range a {
		fs[i] = f(v)
	}
	return fs
}

func MapArrayString[T any](a []T) []string {
	return Fmap[T, string](a, func(v T) string { return fmt.Sprint(v) })
}

func StrJoin(a []string, sep string) string {
	out := ""
	for i, s := range a {
		if i > 0 {
			out += sep
		}
		out += s
	}
	return out
}