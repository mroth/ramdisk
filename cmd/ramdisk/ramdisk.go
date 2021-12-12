// The ramdisk command line application.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mroth/ramdisk"
	"github.com/mroth/ramdisk/datasize"
)

const (
	defaultSize = ramdisk.DefaultSize / datasize.MB
)

func usage() {
	fmt.Fprintf(os.Stderr, `ramdisk 🐏

Usage:
  ramdisk [options] create [<mount-path>]
  ramdisk destroy <device-path>

Options:
  -h -help      Show this screen.
  -v            Verbose output.
  -size=<mb>    Size in megabytes [default: %v].

For more information see https://github.com/mroth/ramdisk
`, defaultSize)
}

func main() {
	size := flag.Uint64("size", defaultSize, "Size in megabytes")
	verbose := flag.Bool("v", false, "Verbose output")
	flag.Usage = usage
	flag.Parse()

	logger := log.New(os.Stderr, "", log.LstdFlags)

	switch flag.Arg(0) {
	case "create":
		var opts ramdisk.Options
		opts.Size = *size * datasize.MB
		if *verbose {
			opts.Logger = logger
		}
		opts.MountPath = flag.Arg(1)

		rd, err := ramdisk.Create(opts)
		if err != nil {
			logger.Println(err)
			// if was ExitError, print full stderr as well
			exiterr, ok := err.(*exec.ExitError)
			if ok {
				logger.Printf("%s\n", exiterr.Stderr)
			}
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "%vMB ramdisk created as %v, mounted at %v\n",
			*size, rd.DevicePath, rd.MountPath)
		fmt.Fprintf(os.Stderr, "To later remove do: `ramdisk destroy %v`\n",
			rd.DevicePath)
	case "destroy":
		device := flag.Arg(1)
		if device == "" {
			flag.Usage()
		}
		err := ramdisk.Destroy(device)
		if err != nil {
			logger.Println(err)
			exiterr, ok := err.(*exec.ExitError)
			if ok {
				logger.Printf("%s\n", exiterr.Stderr)
			}
			os.Exit(1)
		} else if *verbose {
			logger.Printf("Disk %s unmounted\n", device)
		}
	default:
		flag.Usage()
	}
}
