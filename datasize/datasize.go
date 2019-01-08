// Package datasize defines a few useful constants for binary prefix byte sizes.
package datasize

// http://en.wikipedia.org/wiki/Binary_prefix
const (
	B = 1 // byte

	// Decimal
	KB = 1000 * B  // kilobyte
	MB = 1000 * KB // megabyte
	GB = 1000 * MB // gigabyte
	TB = 1000 * GB // terabyte
	PB = 1000 * TB // petabyte
	EB = 1000 * PB // exabyte

	// Binary
	KiB = 1024 * B   // kibibyte
	MiB = 1024 * KiB // mebibyte
	GiB = 1024 * MiB // gibibyte
	TiB = 1024 * GiB // tebibyte
	PiB = 1024 * TiB // pebibyte
	EiB = 1024 * PiB // exbibyte
)
