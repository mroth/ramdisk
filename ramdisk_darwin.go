package ramdisk

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/mroth/ramdisk/datasize"
)

// macosPlatformImplmentation is the implementation for macOS systems.
// https://jakobstoeck.de/2017/ramdisk-for-faster-applications-under-macos/
type macosPlatformImplementation struct{}

func init() {
	implementation = macosPlatformImplementation{}
}

const (
	// new macs seem to have a device block size that is significantly larger,
	// at 4096B, however hdiutil attach apparently still uses old size.
	blockSize uint64 = datasize.B * 512
)

// create creates a macOS ramdisk device, formats it, and mounts it.
//
// Even if an error is generated, a reference to a partial RamDisk struct may be
// returned, as the error may have occured after device creation or formatting
// (for example, a mounting error), and the user may wish to know where the
// orphaned device was left. For example, you could end up with a results such
// as:
//
//   RamDisk{DeviceID:"/dev/disk2", MountPath:""}, Err:"Failed to mount at /foo"
//
func (i macosPlatformImplementation) create(opts Options) (*RamDisk, error) {
	var rd RamDisk

	desiredBlocks := opts.Size / blockSize
	devicePath, deviceErr := createDevice(desiredBlocks)
	if deviceErr != nil {
		return nil, deviceErr
	} else if opts.Logger != nil {
		opts.Logger.Printf("Created ramdisk %s\n", devicePath)
	}
	rd.DevicePath = devicePath

	volumeName := filepath.Base(opts.MountPath)
	formatErr := formatHFS(devicePath, volumeName, opts.Logger)
	if formatErr != nil {
		return &rd, fmt.Errorf("format: %v", formatErr)
	}

	mountErr := mountHFS(devicePath, opts.MountPath, opts.Logger)
	if mountErr != nil {
		return &rd, mountErr
	}

	rd.MountPath = opts.MountPath // set this only once successfully mounted
	return &rd, nil
}

// use hdiutil to create the initial ramdisk device, not yet formatted
// returns the device path (e.g. /dev/disk2) or an error
func createDevice(blockSize uint64) (string, error) {
	path := fmt.Sprintf("ram://%d", blockSize)
	cmd := exec.Command("hdiutil", "attach", "-nomount", path)
	output, err := cmd.Output()
	devicePath := bytes.TrimSpace(output)
	return string(devicePath), err
}

func formatHFS(devicePath, volumeName string, logger *log.Logger) error {
	cmd := exec.Command("newfs_hfs", "-v", volumeName, devicePath)
	stdout, err := cmd.Output()
	if err == nil && logger != nil {
		logger.Printf("%s\n", bytes.TrimSpace(stdout))
	}
	return err
}

func mountHFS(devicePath, mountPath string, logger *log.Logger) error {
	cmd := exec.Command("mount", "-t", "hfs", "-v", devicePath, mountPath)
	stdout, err := cmd.Output()
	if err == nil && logger != nil {
		logger.Printf("%s\n", bytes.TrimSpace(stdout))
	}
	return err
}

func (i macosPlatformImplementation) destroy(devicePath string) error {
	cmd := exec.Command("diskutil", "eject", devicePath)
	return cmd.Run()
}

// unmount apfs: diskutil unmountDisk

// formats and mounts, will end up at /Volumes/<volName>
//
// newfs_apfs will also format, however those dont appear to be mountable via mount_apfs on top of normal filesystem?
//
// https://stackoverflow.com/questions/46224103/create-apfs-ram-disk-on-macos-high-sierra
//
// Usage:  diskutil partitionDisk MountPoint|DiskIdentifier|DeviceNode
//         [numberOfPartitions] [APM[Format]|MBR[Format]|GPT[Format]]
//         [part1Format part1Name part1Size part2Format part2Name part2Size
//          part3Format part3Name part3Size ...]
// diskutil partitionDisk $(hdiutil attach -nomount ram://2048000) 1 GPTFormat APFS 'ramdisk' '100%'
// func formatAPFS(devicePath string, volumeName string) error {
// 	const numPartitions = "1"
// 	const partitionScheme = "GPTFormat"
// 	const partitionFormat = "APFS"
// 	const partitionSize = "100%"
// 	cmd := exec.Command("diskutil",
// 		"partitionDisk", devicePath, numPartitions, partitionScheme,
// 		partitionFormat, volumeName, partitionSize)
// 	return cmd.Run()
// }
