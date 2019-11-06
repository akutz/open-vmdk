package vmdk

import (
	"io"
)

// DiskInfo contains information about a VMDK file.
type DiskInfo interface {
	io.Closer
	io.Seeker
	io.ReaderAt
	io.WriterAt

	Abort() error
	Capacity() int64
}
