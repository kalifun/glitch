package generator

import (
	"github.com/spf13/viper"
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
	c.readFile()
	return nil
}

func (c *codeGenerator) readFile() error {
	viper.SetConfigFile(c.file_path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	var desc ErrorDesc
	err = viper.Unmarshal(&desc)
	if err != nil {
		return err
	}
	return nil
}
