package cmd

import (
	"log"

	"github.com/kalifun/glitch/utils"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:     "gen",
	Short:   "Generate a global error code",
	Long:    "Generate a global error code",
	Example: "glitch gen",
	Run: func(cmd *cobra.Command, args []string) {
		genCode()
	},
}

var yamlFile, packageName string

func init() {
	genCmd.Flags().StringVarP(&packageName, "pack", "p", "errors", "pack name")
	genCmd.Flags().StringVarP(&yamlFile, "yaml", "y", "errors.yaml", "yaml xxx.yaml")
	log.SetFlags(log.Lshortfile)
}

func genCode() {
	exit, _ := utils.PathExists(yamlFile)
	if !exit {
		log.Fatalf("%s file does not exist! use the init command to generate yam files.", yamlFile)
	}
}
