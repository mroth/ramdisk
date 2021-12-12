# ramdisk :ram:

<!-- disable badges until more stable
[![Build Status](https://travis-ci.com/mroth/ramdisk.svg?branch=master)](https://travis-ci.com/mroth/ramdisk)
[![Go Report Card](https://goreportcard.com/badge/github.com/mroth/ramdisk)](https://goreportcard.com/report/github.com/mroth/ramdisk)
[![GoDoc](https://godoc.org/github.com/mroth/ramdisk?status.svg)](https://godoc.org/github.com/mroth/ramdisk)
-->

Convenience wrapper for managing RAM disks across different operating systems
with a consistent interface.

## Usage

Help screen:
```
$ ramdisk -h
ramdisk 0.1.0 🐏

Usage:
  ramdisk [options] create [<mount-path>]
  ramdisk destroy <device-path>

Options:
  -h -help      Show this screen.
  -v            Verbose output.
  -size=<mb>    Size in megabytes [default: 32].
```

Creating a new ram disk with a specified size:
```
$ ramdisk -size=512 create
512MB ramdisk created as /dev/disk5, mounted at /tmp/ramdisk-401987900
To later remove do: `ramdisk destroy /dev/disk5`
```

## Installation

* 💾 Download a [precompiled binary](https://github.com/mroth/ramdisk/releases/).
* 🍺 Homebrew on macOS: `brew install mroth/tap/ramdisk`
* 📦 Compile via Go toolchain: `go install github.com/mroth/ramdisk/cmd/ramdisk@latest`

## Platform Support

### macOS :white_check_mark:

Works great and does not require superuser access.

The basic steps followed are:

- Create an unmounted but attached device in RAM that consists of the
  appropriate number of device blocks via `hdiutil`.
- Format a new uniquely named HFS+ volume on that device via `newfs_hfs`.
- Mounts the volume at a uniquely generated path within the `/tmp` filesystem,
  via `mount`.

This normally requires a sequence of arcane commands on the macOS command line
which I can never remember, and thus was the primary reason I created this
wrapper.

### Linux :white_check_mark:

Things are quite simple and work great via `tmpfs` on Linux! (But note that most
Linux implementations unfortunately requires sudo access to mount new volumes.)

<small>
If you prefer a sudo-less route, most modern Linux on kernel 2.6+ often
_already_ has `/dev/shm` mounted, which is memory backed, so you can also just
use that without any initialization at all.
</small>

### Windows :x:

This would be great, unfortunately there is no built-in support at the operating
system level at current time. The standard solution seems to be to use [ImDisk],
but I haven't built support for that into this library yet because well, relying
on third-party drivers to be installed on an end user system isn't going to
work.

[ImDisk]: http://www.ltr-data.se/opencode.html/#ImDisk