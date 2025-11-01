qm disk import "{{var.vm_id}}" "{{var.cloud_image_filename}}" "{{var.storage_pool}}"
qm set "{{var.vm_id}}" --scsihw virtio-scsi-pci --scsi0 "{{var.storage_pool}}:vm-{{var.vm_id}}-disk-0"
qm resize "{{var.vm_id}}" scsi0 "{{var.disk_size}}G"
qm set "{{var.vm_id}}" --ide2 "{{var.storage_pool}}:cloudinit" --bootdisk scsi0 --serial0 socket --vga serial0
qm set "{{var.vm_id}}" --ipconfig0 ip=dhcp
qm set "{{var.vm_id}}" --ciuser vorpax --sshkeys /root/.ssh/vorpax.pub