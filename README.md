# PDC - Proxmox Direct Config

A sweet CLI tool for managing Proxmox VE config with ease.

PDC automates the creation and management of cloud-init ready VM templates on Proxmox VE hosts via SSH.

## Demo

![demo](https://vhs.charm.sh/vhs-5ykE7QUyezexlKTgWClu3B.gif)

## Features

- **Automated Template Creation**: It downloads cloud-init images and convert them to configurable Proxmox VM templates
- **Remote Management**: No need to run on your node. You can execute commands on remote Proxmox hosts via SSH
- **Configuration**: Save and load VM configurations using YAML or JSON
- **Basic Cloud-Init Support**: Automatic cloud-init configuration with custom user and ssh keys (WIP).
- **A nice and fast experience**: Build with Golang using cobra, Viper, Lipgloss and Goph 

## Prerequisites

- Go 1.25.3 or later (for building from source)
- Access to a Proxmox VE host
- SSH key-based authentication configured for the Proxmox host


## Installation

### From Source

```bash
git clone https://github.com/vorpax/pdc
cd pdc
go build -o pdc
sudo mv pdc /usr/local/bin/
```

## Configuration

PDC supports configuration via YAML files. The tool searches for `config.yaml` in the following locations (in order):

1. `/etc/proxmox-direct-config/`
2. `$HOME/.proxmox-direct-config/`
3. Current directory (`.`)

### Configuration File Example

Create a `config.yaml` file:

```yaml
vm:
  id: "9000"
  name: "ubuntu-template"
  memory: "2048"
  cores: "2"
  disk_size: "32"

proxmox:
  storage: "local-lvm"
  ci_user: "ubuntu"
  ssh_key_file: "/root/.ssh/authorized_keys" # Path on Proxmox host to the public SSH key file 
  bridge: "vmbr0" 

config:
  save: false # Whether you'd like to save the VM config after creation
```

### SSH Configuration

Update the SSH connection settings in `internal/vm.go:connectHostSSH()` to match your environment:

```go
auth, err := goph.Key("/path/to/your/ssh/key", "")
client, err := goph.New("root", "your-proxmox-host", auth)
```

## Usage

### Global Flags

- `-v, --verbose`: Enable verbose output for detailed logging

### Template Create

Create a new Proxmox VM template from a cloud-init image.

```bash
pdc template create [image-url] [flags]
```

#### Flags

- `--vm-id` (required): Unique ID for the new VM/template
- `--vm-name`: Name for the new VM/template
- `--storage` (required): Proxmox storage pool to use
- `--memory`: Memory for the VM in MB (default: "2048")
- `--cores`: Number of CPU cores (default: "2")
- `--disk-size`: Disk size in GB (default: "32")
- `--ci-user`: Cloud-init username
- `--ssh-key-file`: Path on Proxmox host to the public SSH key file
- `--bridge`: Proxmox network bridge (default: "vmbr0")
- `--save`: Save the VM configuration to config.json after creation

#### Example

```bash
# Create an Ubuntu Plucky 25.04 template
pdc template create https://cloud-images.ubuntu.com/plucky/current/plucky-server-cloudimg-amd64.img \
  --vm-id 9000 \
  --vm-name ubuntu-plucky-template \
  --storage local-lvm \
  --memory 2048 \
  --cores 2 \
  --disk-size 32 \
  --ci-user ubuntu \
  --ssh-key-file /root/.ssh/authorized_keys \
  --bridge vmbr0 \
  --save

# Create an Alpine template with verbose output
pdc template create [text](https://dl-cdn.alpinelinux.org/alpine/v3.22/releases/cloud/nocloud_alpine-3.22.2-x86_64-bios-cloudinit-r0.qcow2) \
  --vm-id 9001 \
  --vm-name alpine-3.22-template \
  --storage nfs-storage \
  --verbose
```

### Destroy

```bash
Destroy an existing Proxmox VM template.

```bash
pdc destroy [flags]
```

#### Flags

- `--vm-id`: The unique ID of the VM/template to destroy

#### Example

```bash
# Destroy a VM template
pdc destroy --vm-id 9000

# Destroy with verbose output
pdc destroy --vm-id 9000 --verbose
```

## Template Creation Process

When you run `pdc template create`, the following steps are executed:

1. **Download**: Downloads the cloud-init image to `/tmp` on the Proxmox host
2. **Create VM**: Creates a new VM with specified resources
3. **Import Disk**: Imports the downloaded image as a disk
4. **Attach Disk**: Attaches the disk as scsi0 with virtio-scsi-pci
5. **Resize Disk**: Expands the disk to the specified size
6. **Cloud-Init Setup**: Configures cloud-init drive and boot settings
7. **Network Configuration**: Sets up DHCP networking
8. **User Configuration**: Configures the cloud-init user and SSH keys

## Configuration Storage

When using the `--save` flag, VM configurations are stored in `config.json` in the following format:

```json
{
  "vms": [
    {
      "vm_id": "9000",
      "vm_name": "ubuntu-22.04-template",
      "memory": "2048",
      "cpu_cores": "2",
      "storage_pool": "local-lvm",
      "disk_size": "32",
      "ci_user": "ubuntu",
      "ssh_key_path": "/root/.ssh/authorized_keys",
      "image_url": "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img",
      "network_bridge": "vmbr0"
    }
  ]
}
```


## Troubleshooting

### SSH Connection Issues

If you encounter SSH connection errors:

1. Verify your SSH key path in `internal/vm.go`
2. Ensure the Proxmox host is reachable
3. Confirm SSH key-based authentication is configured
4. Use `--verbose` flag to see detailed error messages

### Storage Pool Errors

If the storage pool is not found:

1. List available storage pools on Proxmox: `pvesm status`
2. Ensure the storage pool name matches exactly
3. Verify you have permissions to use the storage pool

### VM ID Conflicts

If a VM ID already exists:

1. Choose a different ID using `--vm-id`
2. Or destroy the existing VM first using `pdc destroy`

## License

[Your License Here]