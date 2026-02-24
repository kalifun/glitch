package generator

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"mvdan.cc/gofumpt/format"
)

func createOutputDir(outputDir string) error {
	return os.MkdirAll(outputDir, os.ModePerm)
}

func createCodeFile(outputDir, fileName string) (*os.File, error) {
	return os.Create(path.Join(outputDir, fileName))
}

func formatCode(code string) ([]byte, error) {
	return format.Source([]byte(code), format.Options{})
}

func outputFileName(inputPath string, singleInput bool) string {
	if singleInput {
		return "code.go"
	}
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	name = sanitizeFileBase(name)
	return name + ".go"
}

func sanitizeFileBase(name string) string {
	name = strings.ToLower(name)
	re := regexp.MustCompile(`[^a-z0-9_]+`)
	name = re.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	if name == "" {
		name = "errors"
	}
	if name[0] >= '0' && name[0] <= '9' {
		name = "errors_" + name
	}
	return name
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
