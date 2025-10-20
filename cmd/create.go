/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
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

	if err != nil {

	}

	download_template(parsed_url)
}

func download_template(parsed_url *url.URL) (download_path string, err string) {
	filename := parsed_url.Path[strings.LastIndex(parsed_url.Path, "/")+1:]

	// Utiliser le répertoire courant de travail
	wd, wd_err := os.Getwd()

	if wd_err != nil {

		fmt.Println("Error getting working directory:", err)
		return
	}

	download_path = filepath.Join(wd, "data", filename)

	fmt.Printf("Filename: %s\n", filename)
	fmt.Printf("Download path: %s\n", download_path)

	command := exec.Command("curl", "-o", download_path, parsed_url.String())
	fmt.Println("Command:", command)

	if err := command.Run(); err != nil {
		fmt.Println("Error:", err)
	}

	return download_path, err
}
