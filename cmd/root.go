/*
Copyright © 2025 vorpax <git@vorpax.dev>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pdc",
	Short: "Proxmox Direct Config - Manage Proxmox VE configurations with ease",
	Long:  `PDC (Proxmox Direct Config) is a command-line tool for better cluster config. `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.proxmox-direct-config.yaml)")

	// Add a verbose flag
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().String("proxmox-host", "", "Proxmox host address")
	viper.BindPFlag("proxmox.host", rootCmd.PersistentFlags().Lookup("proxmox-host"))

	rootCmd.PersistentFlags().String("proxmox-user", "root", "Proxmox SSH user")
	viper.BindPFlag("proxmox.user", rootCmd.PersistentFlags().Lookup("proxmox-user"))

	rootCmd.PersistentFlags().String("ssh-key", "", "Path to SSH private key for Proxmox connection")
	viper.BindPFlag("proxmox.ssh_key", rootCmd.PersistentFlags().Lookup("ssh-key"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/proxmox-direct-config/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.proxmox-direct-config") // call multiple times to add many search paths
	viper.AddConfigPath(".")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Set environment variable prefix and key replacer
	// This must come AFTER ReadInConfig so flags take precedence
	viper.SetEnvPrefix("PDC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv() // read in environment variables that match
}
