package menv

import "strings"

// Merge merge two env map, if keys in src has a suffix '-', this will delete the key from dst
func Merge(dst map[string]string, src map[string]string) {
	for k, v := range src {
		if strings.HasSuffix(k, "-") {
			delete(dst, k[:len(k)-1])
		} else {
			dst[k] = v
		}
	}
	return
}
