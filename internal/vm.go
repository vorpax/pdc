package internal

import (
	"fmt"
	"log"
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

func CreateTemplate(vmId, vmName, memory, cpuCores, storagePool, diskSize, ciUser, sshKeyPath, imageUrl, networkBridge string, verbose bool) error {
	fmt.Println(InfoStyle.Render("Starting template creation process..."))

	runner, err := createRunner(true) // Always use remote runner for this process
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

func downloadTemplate(parsedUrl *url.URL, execContext CommandRunner) (downloadPath string, err string) {
	filename := filepath.Base(parsedUrl.Path)

	// Utiliser le répertoire courant de travail
	wd, wdErr := os.Getwd()

	if wdErr != nil {

		fmt.Println("Error getting working directory:", err)
		return
	}

	downloadPath = filepath.Join(wd, "data", filename)

	fmt.Printf("Filename: %s\n", filename)
	fmt.Printf("Download path: %s\n", downloadPath)

	args := []string{"-o", downloadPath, parsedUrl.String()}
	cmd, cmdErr := execContext.Command("curl", args...)

	if cmdErr != nil {
		fmt.Println("Error creating command:", cmdErr)
		return downloadPath, cmdErr.Error()
	}

	output, runErr := cmd.CombinedOutput()
	fmt.Println("Command output:", string(output))

	if runErr != nil {
		fmt.Println("Error:", runErr)
		return downloadPath, runErr.Error()
	}

	return downloadPath, ""
}

func createVm(vmId string, memory string, cpuCores string, vmName string, networkBridge string, runner CommandRunner) error {

	vmCreationArgs := []string{"create", vmId, "--memory", memory, "--core", cpuCores, "--name", vmName, "--net0", "virtio"}

	cmd, err := runner.Command("/usr/sbin/qm", vmCreationArgs...)
	if err != nil {
		return fmt.Errorf("error creating vm command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running create vm command: %w. Output: %s", err, string(output))
	}

	fmt.Printf("VM %s created successfully.", vmId)
	return nil
}

func importDisk(runner CommandRunner, vmId string, imagePath string, storagePool string) error {
	args := []string{"disk", "import", vmId, imagePath, storagePool}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return fmt.Errorf("error creating import disk command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running import disk command: %w. Output: %s", err, string(output))
	}

	fmt.Printf("Disk %s imported to VM %s on storage %s.", imagePath, vmId, storagePool)
	return nil
}

func attachDisk(runner CommandRunner, vmId string, storagePool string) error {
	scsiDevice := fmt.Sprintf("%s:vm-%s-disk-0", storagePool, vmId)
	args := []string{"set", vmId, "--scsihw", "virtio-scsi-pci", "--scsi0", scsiDevice}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return fmt.Errorf("error creating attach disk command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running attach disk command: %w. Output: %s", err, string(output))
	}

	fmt.Printf("Disk attached to VM %s as scsi0.", vmId)
	return nil
}

func resizeDisk(runner CommandRunner, vmId string, diskSize string) error {
	size := fmt.Sprintf("%sG", diskSize)
	args := []string{"resize", vmId, "scsi0", size}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return fmt.Errorf("error creating resize disk command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running resize disk command: %w. Output: %s", err, string(output))
	}

	fmt.Printf("Disk resized for VM %s to %s.", vmId, size)
	return nil
}

func configureCloudInit(runner CommandRunner, vmId string, storagePool string) error {
	cloudInitDevice := fmt.Sprintf("%s:cloudinit", storagePool)
	args := []string{"set", vmId, "--ide2", cloudInitDevice, "--bootdisk", "scsi0", "--serial0", "socket", "--vga", "serial0"}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return fmt.Errorf("error creating cloud-init config command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running cloud-init config command: %w. Output: %s", err, string(output))
	}

	fmt.Printf("Cloud-init configured for VM %s.", vmId)
	return nil
}

func configureNetwork(runner CommandRunner, vmId string) error {
	args := []string{"set", vmId, "--ipconfig0", "ip=dhcp"}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return fmt.Errorf("error creating network config command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running network config command: %w. Output: %s", err, string(output))
	}

	fmt.Printf("Network configured for VM %s.", vmId)
	return nil
}

func configureUser(runner CommandRunner, vmId string, user string, sshKeyPath string) error {
	// This assumes the key path is accessible from the Proxmox host.
	args := []string{"set", vmId, "--ciuser", user, "--sshkeys", sshKeyPath}
	cmd, err := runner.Command("/usr/sbin/qm", args...)
	if err != nil {
		return fmt.Errorf("error creating user config command: %w", err)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running user config command: %w. Output: %s", err, string(output))
	}

	fmt.Printf("User and SSH keys configured for VM %s.", vmId)
	return nil
}
func connectHostSSH() (CommandRunner, error) {
	auth, err := goph.Key("/Users/vorpax/.ssh/id_ed25519", "")
	if err != nil {
		return nil, err
	}

	client, err := goph.New("root", "mini-homelab", auth)
	if err != nil {
		log.Printf("%s", err.Error())
		log.Fatal(err)
	}

	return &RemoteRunner{client: client}, nil
}

func createRunner(useRemote bool) (CommandRunner, error) {
	if useRemote {
		return connectHostSSH()
	}
	return &LocalRunner{}, nil
}

func TestCode() {
	// Use LocalRunner for local command execution

	runner, err := createRunner(true)

	if err != nil {
		fmt.Printf("Error connecting to host: %s\n", err)
		return
	}

	cmd, err := runner.Command("hostname")
	if err != nil {
		fmt.Printf("Error creating command: %s\n", err)
		return
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
	} else {
		fmt.Printf("Hostname : %s\n", string(output))
	}

	createVm("1001", "2048", "2", "test-vm", "vmbr0", runner)
}
