/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/url"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var unparsed_url_arg string

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
		createTemplate(args[0])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&unparsed_url_arg, "source", "s", "", "Source directory to read from")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createTemplate(unparsed_url string) {
	parsed_url, err := url.Parse(unparsed_url)

	if unparsed_url == "" {
		fmt.Println("Error: URL cannot be empty")
		return
	}

	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	path := strings.Split(parsed_url.Path, "/")
	filename := path[len(path)-1]
	download_path := strings.Join([]string{"./data/", filename}, "")
	command := exec.Command("curl", "-o", download_path, unparsed_url)
	fmt.Println(command)
	fmt.Println(command)
}
