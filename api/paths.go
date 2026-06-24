package api

import (
	"path"
	"strings"
)

func isPathAllowed(configPaths []string, requestPath string) bool {
	if len(configPaths) == 0 {
		return true
	}

	// Clean the path to prevent directory traversal attacks
	cleanPath := path.Clean("/" + requestPath)
	cleanPath = strings.TrimPrefix(cleanPath, "/")

	if cleanPath == "" || cleanPath == "." {
		return false
	}

	for _, p := range configPaths {
		match, err := path.Match(p, cleanPath)
		if err == nil && match {
			return true
		}
		// Also support checking directory prefixes if pattern doesn't end with wildcard
		if strings.HasPrefix(cleanPath, p+"/") || cleanPath == p {
			return true
		}
	}
	return false
}
