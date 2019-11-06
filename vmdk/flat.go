package vmdk

import "os"

// FlatDiskInfo is information about a flat VMDK file.
type FlatDiskInfo struct {
	DiskInfo
	*os.File
	capacity int64
}

func (d *FlatDiskInfo) Capacity() int64 {
	return d.capacity
}

func (d *FlatDiskInfo) Abort() error {
	return d.File.Close()
}

func FlatOpen(filePath string) (*FlatDiskInfo, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	return &FlatDiskInfo{File: f, capacity: st.Size()}, nil
}

func FlatCreate(filePath string, capacity int64) (*FlatDiskInfo, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	if err := f.Truncate(capacity); err != nil {
		return nil, err
	}
	return &FlatDiskInfo{File: f, capacity: capacity}, nil
}
