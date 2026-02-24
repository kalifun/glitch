package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kalifun/glitch/repo/generator"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:     "gen",
	Short:   "Generate a global error code",
	Long:    "Generate a global error code",
	Example: "glitch gen",
	Run: func(cmd *cobra.Command, args []string) {
		genCode(cmd, args)
	},
}

var yamlFiles []string
var packageName string
var outputDir string

func init() {
	genCmd.Flags().StringVarP(&packageName, "pack", "p", "errors", "pack name")
	genCmd.Flags().StringSliceVarP(&yamlFiles, "yaml", "y", []string{"errors.yaml"}, "yaml file(s), directory, or glob")
	genCmd.Flags().StringVar(&outputDir, "out", "", "output directory (defaults to pack name)")
	log.SetFlags(log.Lshortfile)
}

func genCode(cmd *cobra.Command, args []string) {
	inputs := yamlFiles
	if cmd != nil {
		if !cmd.Flags().Changed("yaml") {
			inputs = nil
		}
	}
	if len(args) > 0 {
		inputs = append(inputs, args...)
	}

	files, err := collectYamlFiles(inputs)
	if err != nil {
		log.Fatalln(err)
	}
	g := generator.NewCodeGen(files, packageName, outputDir)
	if err := g.Exec(); err != nil {
		log.Fatalln(err)
	}
}

func collectYamlFiles(inputs []string) ([]string, error) {
	var out []string
	seen := make(map[string]struct{})

	for _, input := range inputs {
		if input == "" {
			continue
		}

		if hasGlob(input) {
			matches, err := filepath.Glob(input)
			if err != nil {
				return nil, err
			}
			if len(matches) == 0 {
				return nil, fmt.Errorf("no files matched glob: %s", input)
			}
			for _, match := range matches {
				if err := addFileIfYaml(match, seen, &out); err != nil {
					return nil, err
				}
			}
			continue
		}

		info, err := os.Stat(input)
		if err != nil {
			return nil, fmt.Errorf("%s does not exist", input)
		}
		if info.IsDir() {
			yamlFiles, err := collectYamlInDir(input)
			if err != nil {
				return nil, err
			}
			if len(yamlFiles) == 0 {
				return nil, fmt.Errorf("no yaml files found in directory: %s", input)
			}
			for _, file := range yamlFiles {
				if err := addFileIfYaml(file, seen, &out); err != nil {
					return nil, err
				}
			}
			continue
		}

		if err := addFileIfYaml(input, seen, &out); err != nil {
			return nil, err
		}
	}

	sort.Strings(out)
	if len(out) == 0 {
		return nil, fmt.Errorf("no yaml files provided")
	}
	return out, nil
}

func hasGlob(input string) bool {
	return strings.ContainsAny(input, "*?[")
}

func collectYamlInDir(dir string) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if isYamlFile(name) {
			files = append(files, filepath.Join(dir, name))
		}
	}
	return files, nil
}

func isYamlFile(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml")
}

func addFileIfYaml(path string, seen map[string]struct{}, out *[]string) error {
	if !isYamlFile(path) {
		return fmt.Errorf("not a yaml file: %s", path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, ok := seen[abs]; ok {
		return nil
	}
	seen[abs] = struct{}{}
	*out = append(*out, path)
	return nil
}
