package internal

import (
	"fmt"
	"log"
	"net/url"
	"os"
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

func CreateTemplate(unparsedUrl string) {
	if unparsedUrl == "" {
		unparsedUrl = "https://example.com/template.qcow2"
	}

	parsedUrl, err := url.Parse(unparsedUrl)

	if err != nil {
		fmt.Printf("Error parsing URL: %s\n", err)
	} else {
		fmt.Printf("Parsed URL: %s\n", parsedUrl.String())
	}

	downloadTemplate(parsedUrl, nil)
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

	return downloadPath, err
}

func createVm(vmId string, memory string, cpuCores string, vmName string, networkBridge string, runner CommandRunner) {

	vmCreationArgs := []string{"create", vmId, " --memory ", memory, " --core ", cpuCores, "--name", vmName, "--net0 virtio, bridge=", networkBridge}

	cmd, err := runner.Command("/usr/sbin/qm", vmCreationArgs...)

	if err != nil {
		fmt.Println("Error creating command:", err)
		return
	}

	output, err := cmd.Output()
	fmt.Println("Command output:", string(output))

	fmt.Println(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

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

func TestCode() {

	runner, err := connectHostSSH()
	if err != nil {
		fmt.Printf("Error connecting to host: %s\n", err)
		return
	}

	output, err := runner.Run("hostname")

	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
	} else {
		fmt.Printf("Hostname : %s\n", output)
	}

	createVm("1000", "2048", "2", "test-vm", "vmbr0", runner)
}
