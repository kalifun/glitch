package generator

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	Pack       = "package %s"
	Import     = `import "github.com/kalifun/glitch/repo/gerr"`
	Declare    = "// Code generated by glitch. DO NOT EDIT!"
	DeclareErr = `var %s = gerr.ErrWrapper{
	Key:      "%s",
	Code:     "%s",
	Category: "%s",
	Severity: gerr.%s,
	Messages: map[string]string{
%s	},
	Description: "%s",
}`
	RegisterCall = `	if err := gerr.Register(%s); err != nil {
		panic(err)
	}`
)

type codeGenerator struct {
	file_path    string
	package_name string
}

func NewCodeGen(filePaht, packageName string) *codeGenerator {
	return &codeGenerator{
		file_path:    filePaht,
		package_name: packageName,
	}
}

func (c *codeGenerator) Exec() error {
	desc, err := c.readFile()
	if err != nil {
		return err
	}

	err = createPackageDir(c.package_name)
	if err != nil {
		return err
	}

	fileHandler, err := createCodeFile(c.package_name, "code.go")
	if err != nil {
		return err
	}
	defer fileHandler.Close()

	var s string
	s += Declare + "\n"
	s += fmt.Sprintf(Pack, c.package_name) + "\n"
	s += Import + "\n"
	s += desc.ToString()

	data, err := formatCode(s)
	if err != nil {
		return err
	}

	fileHandler.Write(data)
	return nil
}

func (c *codeGenerator) readFile() (ErrorDesc, error) {
	viper.SetConfigFile(c.file_path)
	var desc ErrorDesc
	err := viper.ReadInConfig()
	if err != nil {
		return desc, err
	}
	err = viper.Unmarshal(&desc)
	if err != nil {
		return desc, err
	}
	return desc, nil
}
