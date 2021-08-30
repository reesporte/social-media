package util

import "strings"

// SplitPath splits a path by forward slashes, removing one leading slash.
func SplitPath(path string) []string {
	if path[0] == '/' {
		path = path[1:]
	}
	return strings.Split(path, "/")
}
