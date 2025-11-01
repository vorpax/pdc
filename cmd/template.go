/*
Copyright © 2025 vorpax <git@vorpax.dev>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// templateCmd represents the template command
var TemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage Proxmox VM templates",
	Long: `Manage Proxmox VM templates by creating or modifying cloud-init ready templates.

This command provides subcommands for working with Proxmox VM templates,
including creating new templates from cloud-init images.`,
}

func init() {
	rootCmd.AddCommand(TemplateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
