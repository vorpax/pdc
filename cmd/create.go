// Package cmd pve-eztemplate create /*
package cmd

import (
	"fmt"

	"pve-eztemplate/internal"

	"github.com/spf13/cobra"
)

var unparsedUrlArg string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new template using ez-template",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		//vm.CreateTemplate(args[0])
		internal.TestCode()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&unparsedUrlArg, "source", "s", "", "Source directory to read from")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
