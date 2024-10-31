package main

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func boolToInt(b bool) int {
	switch b {
	case true:  return 1
	default:    return 0
	}
}

func repeatString(s string, n int) string {
	res := ""
	for i := 0; i < n; i++ {
		res += s
	}
	return res
}

// func joinStrings(s []string, sep string) string {
// 	res := ""
// 	for i, v := range s {
// 		if i != 0 {
// 			res += sep
// 		}
// 		res += v
// 	}
// 	return res
// }