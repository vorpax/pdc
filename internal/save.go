package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func loadJsonConfig(configPath string) (map[string]interface{}, error) {

	data, err := os.ReadFile(configPath)

	switch err {
	case os.ErrNotExist:
		os.WriteFile(configPath, []byte(`{"vms": []}`), 0644)
	case nil:
		// File exists, proceed
	default:
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config map[string]interface{}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse JSON config file: %v", err)
	}

	return config, nil

}

func saveJsonConfig(configPath string, config map[string]interface{}) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("FAILED TO MARSHAL JSON CONFIG: %v", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("FAILED TO WRITE JSON CONFIG FILE: %v", err)
	}

	return nil
}

func AddVMToConfig(configPath, vmId, vmName, memory, cpuCores, storagePool, diskSize, ciUser, sshKeyPath, imageUrl, networkBridge string) error {
	config, err := loadJsonConfig(configPath)
	if err != nil {
		return err
	}

	vmConfig := map[string]interface{}{
		"vm_id":          vmId,
		"vm_name":        vmName,
		"memory":         memory,
		"cpu_cores":      cpuCores,
		"storage_pool":   storagePool,
		"disk_size":      diskSize,
		"ci_user":        ciUser,
		"ssh_key_path":   sshKeyPath,
		"image_url":      imageUrl,
		"network_bridge": networkBridge,
	}

	config["vms"] = append(config["vms"].([]interface{}), vmConfig)
	err = saveJsonConfig(configPath, config)
	if err != nil {
		return err
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("VM with ID %s added to config file.", vmId)))
	return nil
}

func SaveTemplateConfig(vmId, vmName, memory, cpuCores, storagePool, diskSize, ciUser, sshKeyPath, imageUrl, networkBridge string, verbose bool) {
	configPath := "config.json"

	err := AddVMToConfig(configPath, vmId, vmName, memory, cpuCores, storagePool, diskSize, ciUser, sshKeyPath, imageUrl, networkBridge)
	if err != nil {
		fmt.Println(ErrorStyle.Render(fmt.Sprintf("Failed to save VM configuration: %v", err)))
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(fmt.Sprintf("VM configuration saved to %s", configPath)))
	}
}
