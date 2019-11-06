package vmdk

import (
	"bytes"
	"fmt"
	"os"
)

// SparseDiskInfo is information about a sparse VMDK file.
type SparseDiskInfo struct {
	DiskInfo
	*os.File
	capacity int64
}

func (d *SparseDiskInfo) Capacity() int64 {
	return d.capacity
}

func (d *SparseDiskInfo) Abort() error {
	return d.File.Close()
}

func SparseOpen(filePath string) (*SparseDiskInfo, error) {
	return nil, nil
}

func CheckSparseExtentHeader(src *SparseExtentHeaderOnDisk) bool {
	return bytes.Equal(src.MagicNumber, leUint32(SPARSE_MAGICNUMBER))
}

func GetSparseExtentHeader(src *SparseExtentHeaderOnDisk) (bool, *SparseExtentHeader) {
	if !CheckSparseExtentHeader(src) {
		return false, nil
	}

	dst := &SparseExtentHeader{}

	dst.Version = cpuUint32(src.Version)
	if dst.Version > SPARSE_VERSION_INCOMPAT_FLAGS {
		return false, nil
	}

	dst.Flags = cpuUint32(src.Flags)
	if bitwiseTrue(
		dst.Flags,
		SPARSEFLAG_INCOMPAT_FLAGS&^SPARSEFLAG_COMPRESSED&^SPARSEFLAG_EMBEDDED_LBA) {
		return false, nil
	}
	if bitwiseTrue(dst.Flags, SPARSEFLAG_VALID_NEWLINE_DETECTOR) {
		if src.SingleEndLineChar != SPARSE_SINGLE_END_LINE_CHAR ||
			src.NonEndLineChar != SPARSE_NON_END_LINE_CHAR ||
			src.DoubleEndLineChar1 != SPARSE_DOUBLE_END_LINE_CHAR1 ||
			src.DoubleEndLineChar2 != SPARSE_DOUBLE_END_LINE_CHAR2 {
			return false, nil
		}
	}
	/* Embedded LBA is allowed with compressed flag only. */
	if bitwiseTrue(dst.Flags, SPARSEFLAG_EMBEDDED_LBA) &&
		!bitwiseTrue(dst.Flags, SPARSEFLAG_COMPRESSED) {
		return false, nil
	}

	dst.CompressAlgorithm = cpuUint16(src.CompressAlgorithm)
	dst.UncleanShutdown = src.UncleanShutdown
	dst.Reserved = 0
	dst.Capacity = SectorType(cpuUint64(src.Capacity))
	dst.GrainSize = SectorType(cpuUint64(src.GrainSize))
	dst.DescriptorOffset = SectorType(cpuUint64(src.DescriptorOffset))
	dst.DescriptorSize = SectorType(cpuUint64(src.DescriptorSize))
	dst.NumGTEsPerGT = cpuUint32(src.NumGTEsPerGT)
	dst.RgdOffset = SectorType(cpuUint64(src.RgdOffset))
	dst.GdOffset = SectorType(cpuUint64(src.GdOffset))
	dst.OverHead = SectorType(cpuUint64(src.OverHead))

	return true, dst
}

func SetSparseExtentHeader(src *SparseExtentHeader, dst *SparseExtentHeaderOnDisk, temp bool) {
	/* Use lowercase 'vmdk' signature for temporary stuff. */
	if temp {
		dst.MagicNumber = leUint32(SPARSE_MAGICNUMBER ^ 0x20202020)
	} else {
		dst.MagicNumber = leUint32(SPARSE_MAGICNUMBER)
	}

	dst.Version = leUint32(src.Version)
	dst.Flags = leUint32(src.Flags)
	dst.SingleEndLineChar = SPARSE_SINGLE_END_LINE_CHAR
	dst.NonEndLineChar = SPARSE_NON_END_LINE_CHAR
	dst.DoubleEndLineChar1 = SPARSE_DOUBLE_END_LINE_CHAR1
	dst.DoubleEndLineChar2 = SPARSE_DOUBLE_END_LINE_CHAR2
	dst.CompressAlgorithm = leUint16(src.CompressAlgorithm)
	dst.UncleanShutdown = src.UncleanShutdown
	dst.Capacity = leUint64(uint64(src.Capacity))
	dst.GrainSize = leUint64(uint64(src.GrainSize))
	dst.DescriptorOffset = leUint64(uint64(src.DescriptorOffset))
	dst.DescriptorSize = leUint64(uint64(src.DescriptorSize))
	dst.NumGTEsPerGT = leUint32(src.NumGTEsPerGT)
	dst.RgdOffset = leUint64(uint64(src.RgdOffset))
	dst.GdOffset = leUint64(uint64(src.GdOffset))
	dst.OverHead = leUint64(uint64(src.OverHead))
}

const diskDescriptorFileFormat = `# Disk DescriptorFile
version=1
encoding="UTF-8"
CID=%08x
parentCID=ffffffff
createType="streamOptimized"

# Extent description
RW %d SPARSE "%s"

# The Disk Data Base
#DDB

ddb.longContentID = "%08x%08x%08x%08x"
ddb.toolsVersion = "2147483647"
/* OpenSource Tools version. */
ddb.virtualHWVersion = "4"
/* This field is obsolete, used by ESX3.x and older only. */
ddb.geometry.cylinders = "%d"
ddb.geometry.heads = "255"
/* 255/63 is good for anything bigger than 4GB. */
ddb.geometry.sectors = "63"
ddb.adapterType = "lsilogic"
`

func MakeDiskDescriptorFile(fileName string, capacity uint64, contentID uint32) string {

	var cylinders uint64

	if capacity > 65535*255*63 {
		cylinders = 65535
	} else {
		cylinders = ceiling(capacity, 255*63)
	}

	return fmt.Sprintf(
		diskDescriptorFileFormat,
		contentID,
		capacity,
		fileName,
		mrand48.Uint32(),
		mrand48.Uint32(),
		mrand48.Uint32(),
		contentID,
		cylinders)
}

type ZLibBuffer struct {
	grainHdr   *SparseGrainLBAHeaderOnDisk
	specialHdr *SparseSpecialLBAHeaderOnDisk
	data       []byte
}

type SparseGTInfo struct {
	GTEs          uint64
	GTs           uint32
	GDsectors     uint32
	GTsectors     uint32
	LastGrainNr   uint64
	LastGrainSize uint32
	Gd            *LEUint32
	Gt            *LEUint32
}

type SparseVmdkWriter struct {
	GtInfo     SparseGTInfo
	GdOffset   int64
	GtOffset   int64
	RgdOffset  int64
	RgtOffset  int64
	CurSP      uint32
	ZLibBuffer ZLibBuffer
	//z_stream     zstream
	//int          fd
	FileName    string
	GrainBuffer []byte
	//uint64_t     grainBufferNr
	//uint32_t     grainBufferValidStart
	//uint32_t     grainBufferValidEnd
}
