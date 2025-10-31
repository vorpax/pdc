package internal

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/melbahja/goph"
)

type CommandRunner interface {
	Run(command string, args ...string) (string, error)
}

type LocalRunner struct{}

func (lr *LocalRunner) Run(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

type RemoteRunner struct {
	client *goph.Client
}

func (rr *RemoteRunner) Run(command string, args ...string) (string, error) {
	fullCommand := command
	if len(args) > 0 {
		fullCommand += " " + strings.Join(args, " ")
	}
	output, err := rr.client.Run(fullCommand)
	return string(output), err
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
	filename := parsedUrl.Path[strings.LastIndex(parsedUrl.Path, "/")+1:]

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
	command, runErr := execContext.Run("curl", args...)
	fmt.Println("Command:", command)

	if runErr != nil {
		fmt.Println("Error:", err)
	}

	return downloadPath, err
}

func createVm(vmId string, memory string, cpuCores string, vmName string, networkBridge string, runner CommandRunner) {

	vmCreationArgs := []string{"create", vmId, " --memory ", memory, " --core ", cpuCores, "--name", vmName, "--net0 virtio, bridge=", networkBridge}

	command, err := runner.Run("/usr/sbin/qm", vmCreationArgs...)

	fmt.Println("Command:", command)

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
		log.Printf(err.Error())
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
