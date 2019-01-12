package ramdisk

import (
	"bytes"
	"fmt"
	"os/exec"
)

// LinuxPlatformImplementation is the implementation for Linux systems.
//
// The Linux implementation likely *requires sudo* to function on most distros.
// If you want a sudo-less option, you can simply use /dev/shm instead on most
// modern Linux platforms.
type LinuxPlatformImplementation struct{}

func init() {
	implementation = LinuxPlatformImplementation{}
}

func (i LinuxPlatformImplementation) create(opts Options) (*RAMDisk, error) {
	rd := RAMDisk{DevicePath: opts.MountPath, MountPath: opts.MountPath}
	sizeFlag := fmt.Sprintf("size=%d", opts.Size)
	cmd := exec.Command(
		"mount", "-v", "-t", "tmpfs", "-o", sizeFlag, "tmpfs", opts.MountPath)
	stdout, err := cmd.Output()
	if err == nil && opts.Logger != nil {
		opts.Logger.Printf("%s\n", bytes.TrimSpace(stdout))
	}
	return &rd, err
}

func (i LinuxPlatformImplementation) destroy(devicePath string) error {
	cmd := exec.Command("umount", devicePath)
	return cmd.Run()
}
