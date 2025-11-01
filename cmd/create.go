package cmd

import (
	"fmt"
	"log"

	"pve-eztemplate/internal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [image-url]",
	Short: "Create a new Proxmox template from a cloud-init image",
	Long: `This command automates the creation of a Proxmox VM template.
It downloads a specified cloud-init ready image, creates a new VM,
configures it with the provided settings (or settings from a config file),
and converts it into a template for future use.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		imageUrl := args[0]

		err := internal.CreateTemplate(
			viper.GetString("vm.id"),
			viper.GetString("vm.name"),
			viper.GetString("vm.memory"),
			viper.GetString("vm.cores"),
			viper.GetString("proxmox.storage"),
			viper.GetString("vm.disk_size"),
			viper.GetString("proxmox.ci_user"),
			viper.GetString("proxmox.ssh_key_file"),
			imageUrl,
			viper.GetString("proxmox.bridge"),
			viper.GetBool("verbose"),
		)

		if viper.GetBool("config.save") {
			internal.SaveTemplateConfig(
				viper.GetString("vm.id"),
				viper.GetString("vm.name"),
				viper.GetString("vm.memory"),
				viper.GetString("vm.cores"),
				viper.GetString("proxmox.storage"),
				viper.GetString("vm.disk_size"),
				viper.GetString("proxmox.ci_user"),
				viper.GetString("proxmox.ssh_key_file"),
				imageUrl,
				viper.GetString("proxmox.bridge"),
				viper.GetBool("verbose"),
			)

		}

		if err != nil {
			log.Fatalf("%s", internal.ErrorStyle.Render(fmt.Sprintf("Template creation failed: %v", err)))
		}
	},
}

func init() {

	TemplateCmd.AddCommand(createCmd)
	// Define flags and bind them to Viper
	createCmd.Flags().String("vm-id", "", "The unique ID for the new VM/template")
	viper.BindPFlag("vm.id", createCmd.Flags().Lookup("vm-id"))

	createCmd.Flags().String("vm-name", "", "Name for the new VM/template")
	viper.BindPFlag("vm.name", createCmd.Flags().Lookup("vm-name"))

	createCmd.Flags().String("memory", "2048", "Memory for the VM in MB")
	viper.BindPFlag("vm.memory", createCmd.Flags().Lookup("memory"))

	createCmd.Flags().String("cores", "2", "Number of CPU cores for the VM")
	viper.BindPFlag("vm.cores", createCmd.Flags().Lookup("cores"))

	createCmd.Flags().String("storage", "", "Proxmox storage pool to use")
	viper.BindPFlag("proxmox.storage", createCmd.Flags().Lookup("storage"))

	createCmd.Flags().String("disk-size", "32", "Disk size for the VM in GB")
	viper.BindPFlag("vm.disk_size", createCmd.Flags().Lookup("disk-size"))

	createCmd.Flags().String("ci-user", "", "Cloud-init username")
	viper.BindPFlag("proxmox.ci_user", createCmd.Flags().Lookup("ci-user"))

	createCmd.Flags().String("ssh-key-file", "", "Path on the Proxmox host to the public SSH key file")
	viper.BindPFlag("proxmox.ssh_key_file", createCmd.Flags().Lookup("ssh-key-file"))

	createCmd.Flags().String("bridge", "vmbr0", "Proxmox network bridge")
	viper.BindPFlag("proxmox.bridge", createCmd.Flags().Lookup("bridge"))

	createCmd.Flags().Bool("save", false, "Save the VM configuration to the config file after creation")
	viper.BindPFlag("config.save", createCmd.Flags().Lookup("save"))

	// Mark essential flags as required
	createCmd.MarkFlagRequired("vm-id")
	createCmd.MarkFlagRequired("storage")
}
