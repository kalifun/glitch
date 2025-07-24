package cmd

import (
	"os"
	"path"

	"github.com/kalifun/glitch/utils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "init error template",
	Long:    "init error template",
	Example: "glitch init",
	Run: func(cmd *cobra.Command, args []string) {
		init_template()
	},
}

const errorTmp = `error:
  - key: MissingParameter
    code: MissingParameter
	category: validation
    severity: error
    description: Parameter verification error
    message:
      cn: "缺少参数: %s"
      en: "Missing Parameter: %s"`

func init_template() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filePath := path.Join(dir, "errors.yaml")
	ok, err := utils.PathExists(filePath)
	if err != nil {
		panic(err)
	}

	if !ok {
		utils.WiteFile(filePath, errorTmp)
	}
}
