
Output demo.gif

Set Shell "zsh"
Set FontSize 14
Set Width 1200
Set Height 600
Set Padding 20
Set TypingSpeed 50ms
Set Theme catppuccin-mocha
Set WindowBar Colorful

# Show the help command
Type "./pdc --help"
Enter
Sleep 2s

# Clear and show template help
Type "clear"
Enter
Sleep 500ms
Type "./pdc template create --help"
Enter
Sleep 3s

# Clear screen
Type "clear"
Enter
Sleep 500ms

# Show a create command example (not actually running it)
Type "# Example: Create an Ubuntu 24.04 template"
Enter
Sleep 500ms
Type "./pdc template create https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img \\"
Enter
Type "  --vm-id 9000 \\"
Enter
Type "  --vm-name ubuntu-24.04-template \\"
Enter
Type "  --storage local-lvm \\"
Enter
Type "  --proxmox-host proxmox.local \\"
Enter
Type "  --ssh-key ~/.ssh/id_ed25519 \\"
Enter
Type "  --ci-user ubuntu \\"
Enter
Type "  --ssh-key-file /root/.ssh/authorized_keys \\"
Enter
Type "  --verbose"
Enter

Sleep 2s

# Show config file example
Type "clear"
Enter
Sleep 500ms
Type "# Or use a config file:"
Enter
Sleep 1s
Type "cat config.yaml"
Enter
Sleep 500ms
Type "proxmox:"
Enter
Type "  host: proxmox.local"
Enter
Type "  user: root"
Enter
Type "  storage: local-lvm"
Enter
Type "  ci_user: ubuntu"
Enter
Type ""
Enter
Type "vm:"
Enter
Type "  id: 9000"
Enter
Type "  name: ubuntu-24.04-template"
Enter
Sleep 2s

Type "clear"
Enter
Sleep 500ms
Type "# Then run with simpler command:"
Enter
Sleep 500ms
Type "./pdc template create https://cloud-images.ubuntu.com/.../image.img"
Enter

Sleep 3s
