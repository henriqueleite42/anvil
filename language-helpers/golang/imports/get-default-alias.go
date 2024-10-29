package imports

import "strings"

func GetDefaultAlias(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
