/*
Copyright © 2025 vorpax <git@vorpax.dev>
*/
package cmd

import (
	"github.com/vorpax/pdc/internal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an existing Proxmox VM template",
	Long: `Destroy an existing Proxmox VM or template by its ID.

This command will permanently remove the specified VM/template from your Proxmox host.
Use with caution as this operation cannot be undone.

Example:
  pdc destroy --vm-id 9000`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		destroyTemplate(
			viper.GetString("vm.id"),
			viper.GetString("proxmox.host"),
			viper.GetString("proxmox.user"),
			viper.GetString("proxmox.ssh_key"),
			viper.GetBool("verbose"),
		)
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().String("vm-id", "", "The unique ID of the VM/template to destroy")
	viper.BindPFlag("vm.id", destroyCmd.Flags().Lookup("vm-id"))

	// Mark vm-id as required
	destroyCmd.MarkFlagRequired("vm-id")
}

func destroyTemplate(vmId, proxmoxHost, proxmoxUser, proxmoxSSHKey string, verbose bool) error {
	internal.DestroyVm(vmId, proxmoxHost, proxmoxUser, proxmoxSSHKey, verbose)

	return nil
}
