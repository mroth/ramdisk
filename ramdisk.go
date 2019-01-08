package ramdisk

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/mroth/ramdisk/datasize"
)

const Version = "0.0.x"

// defaults that are used for any zero value in Options
const (
	DefaultSize = 32 * datasize.MiB
)

// Options are optional values that will override default behavior
type Options struct {
	MountPath string      // optional mountpath (if zero, fileutil.TmpDir will be used)
	Size      uint64      // optional size in bytes (if zero, constant DefaultSize will be used)
	Logger    *log.Logger // if defined, steps may be more verbosely logged to it
}

// RamDisk represents the "results" of a ram disk creation operation
type RamDisk struct {
	// The system path referring to the ramdisk. This may or may not be
	// identical to the MountPath, depending on operating system specific
	// implementations.
	DevicePath string
	// The filesystem path where the ramdisk is mounted and may be viewed.
	MountPath string
}

type PlatformImplementation interface {
	create(opts Options) (*RamDisk, error)
	destroy(deviceID string) error
}

// should be assigned via build constraint'd pkg
var implementation PlatformImplementation

// Create a new ramdisk, using the implementation for the currently active
// platform.
//
// If you wish to use all default values, simply supply a zero-value struct.
//
//     rd, err := ramdisk.Create(Options{})
//
// May return an error on numerous conditions.
func Create(opts Options) (*RamDisk, error) {
	if implementation == nil {
		return nil, errors.New("platform not supported")
	}
	if err := opts.applyDefaults(); err != nil {
		return nil, err
	}

	return implementation.create(opts)
}

func (o *Options) applyDefaults() error {
	if o.Size == 0 {
		o.Size = DefaultSize
	}
	if o.MountPath == "" {
		tmp, err := ioutil.TempDir("", "ramdisk-")
		if err != nil {
			return err
		}
		o.MountPath = tmp
	}
	return nil
}

// Destroy unmounts the ramdisk and removes it from devices.
func Destroy(devicePath string) error {
	if implementation == nil {
		return errors.New("platform not supported")
	}
	return implementation.destroy(devicePath)
}
