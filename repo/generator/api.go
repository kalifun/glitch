package generator

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
	"mvdan.cc/gofumpt/format"
)

const (
	Pack       = "package %s"
	Import     = `import "github.com/kalifun/glitch/repo/gerr"`
	Declare    = "// Code generated by glitch. DO NOT EDIT!"
	DeclareErr = `var %s = gerr.ErrorCode{
	Key:  "%s",
	Code: "%s",
	Message: gerr.ErrorMessage{
		CN: "%s",
		EN: "%s",
	},
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

	err = os.MkdirAll(c.package_name, os.ModePerm)
	if err != nil {
		return err
	}

	fileHandler, err := os.Create(path.Join(c.package_name, "code.go"))
	if nil != err {
		return err
	}
	defer fileHandler.Close()
	var s string
	s += Declare + "\n"
	s += fmt.Sprintf(Pack, c.package_name) + "\n"
	s += Import + "\n"
	s += desc.ToString()

	data, err := format.Source([]byte(s), format.Options{})
	if err != nil {
		log.Fatalf("Error: %s\n%s", err.Error(), s)
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
