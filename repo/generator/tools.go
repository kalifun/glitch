package generator

import (
	"os"
	"path"
	"strings"

	"mvdan.cc/gofumpt/format"
)

func createPackageDir(packageName string) error {
	return os.MkdirAll(packageName, os.ModePerm)
}

func createCodeFile(packageName, fileName string) (*os.File, error) {
	return os.Create(path.Join(packageName, fileName))
}

func formatCode(code string) ([]byte, error) {
	return format.Source([]byte(code), format.Options{})
}

// auto check if string needs to be formatted
// support %v, %s, %d, %f, %t, %b, %c, %o, %x, %X, %U, %e, %E, %f, %F, %g, %G, %p, %q, %x, %X, %U.
func needsFormat(s string) bool {
	// Skip the escaped %%
	s = strings.ReplaceAll(s, "%%", "")

	// Check the remaining %
	for i := 0; i < len(s); i++ {
		if s[i] == '%' {
			if i+1 >= len(s) {
				return false // Single % at the end
			}
			next := s[i+1]
			if strings.ContainsRune("vTtbcdoqxXUeEfFgGsp", rune(next)) {
				return true
			}
		}
	}
	return false
}
