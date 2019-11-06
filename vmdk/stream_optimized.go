package vmdk

import "os"

// StreamOptimizedDiskInfo is information about a stream-optimized VMDK file.
type StreamOptimizedDiskInfo struct {
	DiskInfo
	*os.File
	capacity int64
}

func (d *StreamOptimizedDiskInfo) Capacity() int64 {
	return d.capacity
}

func (d *StreamOptimizedDiskInfo) Abort() error {
	return d.File.Close()
}

func StreamOptimizedCreate(
	filePath string, capacity int) (*StreamOptimizedDiskInfo, error) {

	return nil, nil
}
