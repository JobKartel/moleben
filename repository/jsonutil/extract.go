package jsonutil

func ExtractJSONObject(s string) (string, bool) {
	start := -1
	depth := 0
	for i, r := range s {
		if r == '{' {
			if depth == 0 {
				start = i
			}
			depth++
		}
		if r == '}' {
			if depth > 0 {
				depth--
				if depth == 0 && start != -1 {
					return s[start : i+1], true
				}
			}
		}
	}
	return "", false
}
