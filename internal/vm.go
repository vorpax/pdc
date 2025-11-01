package internal

import (
	"fmt"

	"net/url"
	"os/exec"
	"path/filepath"

	"github.com/melbahja/goph"
)

// Cmd interface represents both os/exec.Cmd and goph.Cmd
type Cmd interface {
	Run() error
	CombinedOutput() ([]byte, error)
	Output() ([]byte, error)
	Start() error
	Wait() error
	// String() ([]byte, error)
}

type CommandRunner interface {
	Command(command string, args ...string) (Cmd, error)
}

type LocalRunner struct{}

func (lr *LocalRunner) Command(command string, args ...string) (Cmd, error) {
	return exec.Command(command, args...), nil
}

type RemoteRunner struct {
	client *goph.Client
}

func (rr *RemoteRunner) Command(command string, args ...string) (Cmd, error) {
	return rr.client.Command(command, args...)
}

func CreateTemplate(vmId, vmName, memory, cpuCores, storagePool, diskSize, ciUser, sshKeyPath, imageUrl, networkBridge, proxmoxHost, proxmoxUser, proxmoxSSHKey string, verbose bool) error {
	fmt.Println(InfoStyle.Render("Starting template creation process..."))

	runner, err := createRunner(true, proxmoxHost, proxmoxUser, proxmoxSSHKey) // Always use remote runner for this process
	if err != nil {
		return fmt.Errorf("failed to connect to Proxmox host: %w", err)
	}

	parsedUrl, err := url.Parse(imageUrl)
	if err != nil {
		return fmt.Errorf("invalid image URL: %w", err)
	}

	fmt.Printf("Downloading image from %s on Proxmox host...\n", imageUrl)
	imagePath, output, err := downloadTemplate(parsedUrl, runner, verbose)
	if err != nil {
		return fmt.Errorf("failed to download template image: %w", err)
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(output))
	}

	output, err = createVm(vmId, memory, cpuCores, vmName, networkBridge, runner, verbose)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(output))
	}

	output, err = importDisk(runner, vmId, imagePath, storagePool, verbose)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(output))
	}

	output, err = attachDisk(runner, vmId, storagePool, verbose)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(output))
	}

	output, err = resizeDisk(runner, vmId, diskSize, verbose)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(output))
	}

	output, err = configureCloudInit(runner, vmId, storagePool, verbose)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(output))
	}

	output, err = configureNetwork(runner, vmId, verbose)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(output))
	}

	output, err = configureUser(runner, vmId, ciUser, sshKeyPath, verbose)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(VerboseStyle.Render(output))
	}

	fmt.Println(SuccessStyle.Render("Template creation process completed successfully!"))
	return nil
}

func downloadTemplate(parsedUrl *url.URL, execContext CommandRunner, verbose bool) (string, string, error) {
	filename := filepath.Base(parsedUrl.Path)
	downloadPath := filepath.Join("/tmp", filename)

	fmt.Println(InfoStyle.Render(fmt.Sprintf("Image will be downloaded to: %s on the remote host", downloadPath)))

	args := []string{"-o", downloadPath, parsedUrl.String()}
	cmd, err := execContext.Command("curl", args...)
	if err != nil {
		return "", "", fmt.Errorf("error creating download command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", string(output), fmt.Errorf("error running download command: %w", err)
	}

	return downloadPath, string(output), nil
}

func createVm(vmId string, memory string, cpuCores string, vmName string, networkBridge string, runner CommandRunner, verbose bool) (string, error) {
	vmCreationArgs := []string{"create", vmId, "--memory", memory, "--core", cpuCores, "--name", vmName, "--net0", "virtio"}

	cmd, err := runner.Command("/usr/sbin/qm", vmCreationArgs...)
	if err != nil {
		return "", fmt.Errorf("error creating vm command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error running create vm command: %w", err)
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("VM %s created successfully.", vmId)))
	return string(output), nil
}

func importDisk(runner CommandRunner, vmId string, imagePath string, storagePool string, verbose bool) (string, error) {
	args := []string{"disk", "import", vmId, imagePath, storagePool}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return "", fmt.Errorf("error creating import disk command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error running import disk command: %w", err)
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("Disk %s imported to VM %s on storage %s.", imagePath, vmId, storagePool)))
	return string(output), nil
}

func attachDisk(runner CommandRunner, vmId string, storagePool string, verbose bool) (string, error) {
	scsiDevice := fmt.Sprintf("%s:vm-%s-disk-0", storagePool, vmId)
	args := []string{"set", vmId, "--scsihw", "virtio-scsi-pci", "--scsi0", scsiDevice}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return "", fmt.Errorf("error creating attach disk command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error running attach disk command: %w", err)
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("Disk attached to VM %s as scsi0.", vmId)))
	return string(output), nil
}

func resizeDisk(runner CommandRunner, vmId string, diskSize string, verbose bool) (string, error) {
	size := fmt.Sprintf("%sG", diskSize)
	args := []string{"resize", vmId, "scsi0", size}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return "", fmt.Errorf("error creating resize disk command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error running resize disk command: %w", err)
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("Disk resized for VM %s to %s.", vmId, size)))
	return string(output), nil
}

func configureCloudInit(runner CommandRunner, vmId string, storagePool string, verbose bool) (string, error) {
	cloudInitDevice := fmt.Sprintf("%s:cloudinit", storagePool)
	args := []string{"set", vmId, "--ide2", cloudInitDevice, "--bootdisk", "scsi0", "--serial0", "socket", "--vga", "serial0"}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return "", fmt.Errorf("error creating cloud-init config command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error running cloud-init config command: %w", err)
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("Cloud-init configured for VM %s.", vmId)))
	return string(output), nil
}

func configureNetwork(runner CommandRunner, vmId string, verbose bool) (string, error) {
	args := []string{"set", vmId, "--ipconfig0", "ip=dhcp"}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return "", fmt.Errorf("error creating network config command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error running network config command: %w", err)
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("Network configured for VM %s.", vmId)))
	return string(output), nil
}

func configureUser(runner CommandRunner, vmId string, user string, sshKeyPath string, verbose bool) (string, error) {
	// This assumes the key path is accessible from the Proxmox host.
	args := []string{"set", vmId, "--ciuser", user, "--sshkeys", sshKeyPath}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return "", fmt.Errorf("error creating user config command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error running user config command: %w", err)
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("User and SSH keys configured for VM %s.", vmId)))
	return string(output), nil
}
func connectHostSSH(host, user, sshKeyPath string) (CommandRunner, error) {
	if host == "" {
		return nil, fmt.Errorf("proxmox host is required (use --proxmox-host flag or set proxmox.host in config)")
	}
	if sshKeyPath == "" {
		return nil, fmt.Errorf("ssh key path is required (use --ssh-key flag or set proxmox.ssh_key in config)")
	}

	auth, err := goph.Key(sshKeyPath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load SSH key from %s: %w", sshKeyPath, err)
	}

	client, err := goph.New(user, host, auth)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s@%s: %w", user, host, err)
	}

	return &RemoteRunner{client: client}, nil
}

func createRunner(useRemote bool, host, user, sshKeyPath string) (CommandRunner, error) {
	if useRemote {
		return connectHostSSH(host, user, sshKeyPath)
	}
	return &LocalRunner{}, nil
}

func DestroyVm(vmId, proxmoxHost, proxmoxUser, proxmoxSSHKey string, verbose bool) (string, error) {

	runner, err := createRunner(true, proxmoxHost, proxmoxUser, proxmoxSSHKey) // Always use remote runner for this process
	if err != nil {
		return "", fmt.Errorf("failed to connect to Proxmox host: %w", err)
	}

	args := []string{"destroy", vmId}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return "", fmt.Errorf("error creating destroy command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error running destroy command: %w", err)
	}

	fmt.Println(VerboseStyle.Render(fmt.Sprintf("Destroying VM %s...", vmId)))
	fmt.Println(VerboseStyle.Render(fmt.Sprintf("destroy command output: *", output)))

	fmt.Println(InfoStyle.Render(fmt.Sprintf("VM %s destroyed.", vmId)))
	return string(output), nil
}
