package ramdisk

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/mroth/ramdisk/datasize"
)

// Version is the semantic version of this package.
const Version = "0.1.0"

// defaults that are used for any zero value in Options
const (
	DefaultSize = 32 * datasize.MB
)

// Options are optional values that will override default behavior
type Options struct {
	MountPath string      // optional: fs mount dir  (default: temp directory)
	Size      uint64      // optional: size in bytes (default: DefaultSize)
	Logger    *log.Logger // optional: logger for verbose implementation logs
}

// RAMDisk represents the "results" of a ram disk creation operation
type RAMDisk struct {
	// The system path referring to the RAMDisk. This may or may not be
	// identical to the MountPath, depending on operating system specific
	// implementations.
	DevicePath string
	// The filesystem path where the RAMDisk is mounted and may be viewed.
	MountPath string
}

// PlatformImplementation is the interface that should be implmented on an
// individual platform (operating system, etc), and hidden behind platform
// specific build tags.
type PlatformImplementation interface {
	create(opts Options) (*RAMDisk, error)
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
// May return an error on numerous platform-specific conditions.
func Create(opts Options) (*RAMDisk, error) {
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
		tmpdir := os.TempDir()
		// the default TempDir() on darwin is designed to be unpredictable and
		// secure, however that makes its a long ugly monstrosity that makes it
		// terrible from a UX perspective if presented to the end-user directly.
		if runtime.GOOS == "darwin" {
			tmpdir = "/tmp"
		}
		tmp, err := ioutil.TempDir(tmpdir, "ramdisk-")
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
