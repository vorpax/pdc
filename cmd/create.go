package cmd

import (
	"fmt"
	"log"

	"github.com/vorpax/pdc/internal"

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

		// Helper to get value from flag (flat key) or config (nested key)
		getConfig := func(flagKey, configKey string) string {
			if val := viper.GetString(flagKey); val != "" {
				return val
			}
			return viper.GetString(configKey)
		}

		vmId := getConfig("vm-id", "vm.id")
		vmName := getConfig("vm-name", "vm.name")
		memory := getConfig("memory", "vm.memory")
		cores := getConfig("cores", "vm.cores")
		storage := getConfig("storage", "proxmox.storage")
		diskSize := getConfig("disk-size", "vm.disk_size")
		ciUser := getConfig("ci-user", "proxmox.ci_user")
		sshKeyFile := getConfig("ssh-key-file", "proxmox.ssh_key_file")
		bridge := getConfig("bridge", "proxmox.bridge")
		save := viper.GetBool("save") || viper.GetBool("config.save")

		err := internal.CreateTemplate(
			vmId,
			vmName,
			memory,
			cores,
			storage,
			diskSize,
			ciUser,
			sshKeyFile,
			imageUrl,
			bridge,
			viper.GetString("proxmox.host"),
			viper.GetString("proxmox.user"),
			viper.GetString("proxmox.ssh_key"),
			viper.GetBool("verbose"),
		)

		if save {
			internal.SaveTemplateConfig(
				vmId,
				vmName,
				memory,
				cores,
				storage,
				diskSize,
				ciUser,
				sshKeyFile,
				imageUrl,
				bridge,
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

	// Define flags
	createCmd.Flags().String("vm-id", "", "The unique ID for the new VM/template")
	createCmd.Flags().String("vm-name", "", "Name for the new VM/template")
	createCmd.Flags().String("memory", "2048", "Memory for the VM in MB")
	createCmd.Flags().String("cores", "2", "Number of CPU cores for the VM")
	createCmd.Flags().String("storage", "", "Proxmox storage pool to use")
	createCmd.Flags().String("disk-size", "32", "Disk size for the VM in GB")
	createCmd.Flags().String("ci-user", "", "Cloud-init username")
	createCmd.Flags().String("ssh-key-file", "", "Path on the Proxmox host to the public SSH key file")
	createCmd.Flags().String("bridge", "vmbr0", "Proxmox network bridge")
	createCmd.Flags().Bool("save", false, "Save the VM configuration to the config file after creation")

	// Bind flags to Viper - use flat keys to avoid nested key issues with config file
	viper.BindPFlag("vm-id", createCmd.Flags().Lookup("vm-id"))
	viper.BindPFlag("vm-name", createCmd.Flags().Lookup("vm-name"))
	viper.BindPFlag("memory", createCmd.Flags().Lookup("memory"))
	viper.BindPFlag("cores", createCmd.Flags().Lookup("cores"))
	viper.BindPFlag("storage", createCmd.Flags().Lookup("storage"))
	viper.BindPFlag("disk-size", createCmd.Flags().Lookup("disk-size"))
	viper.BindPFlag("ci-user", createCmd.Flags().Lookup("ci-user"))
	viper.BindPFlag("ssh-key-file", createCmd.Flags().Lookup("ssh-key-file"))
	viper.BindPFlag("bridge", createCmd.Flags().Lookup("bridge"))
	viper.BindPFlag("save", createCmd.Flags().Lookup("save"))

	// Mark essential flags as required
	createCmd.MarkFlagRequired("vm-id")
	createCmd.MarkFlagRequired("storage")
}
