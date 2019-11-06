package vmdk

/* ********************************************************************************
 * Copyright (c) 2014 VMware, Inc.  All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the “License”); you may not
 * use this file except in compliance with the License.  You may obtain a copy of
 * the License at:
 *
 *            http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an “AS IS” BASIS, without warranties or
 * conditions of any kind, EITHER EXPRESS OR IMPLIED.  See the License for the
 * specific language governing permissions and limitations under the License.
 * *********************************************************************************/

import (
	"encoding/binary"
	"math/rand"
	"time"
)

const (
	SPARSE_MAGICNUMBER                    = 0x564d444b /* VMDK */
	SPARSE_VERSION_INCOMPAT_FLAGS         = 3
	SPARSE_GTE_EMPTY                      = 0x00000000
	SPARSE_GD_AT_END                 uint = 0xFFFFFFFFFFFFFFFF
	SPARSE_SINGLE_END_LINE_CHAR      byte = '\n'
	SPARSE_NON_END_LINE_CHAR         byte = ' '
	SPARSE_DOUBLE_END_LINE_CHAR1     byte = '\r'
	SPARSE_DOUBLE_END_LINE_CHAR2     byte = '\n'
	SPARSE_COMPRESSALGORITHM_NONE         = 0x0000
	SPARSE_COMPRESSALGORITHM_DEFLATE      = 0x0001
)

const (
	SPARSEFLAG_COMPAT_FLAGS           uint32 = 0x0000FFFF
	SPARSEFLAG_VALID_NEWLINE_DETECTOR        = (1 << 0)
	SPARSEFLAG_USE_REDUNDANT                 = (1 << 1)
	SPARSEFLAG_MAGIC_GTE                     = (1 << 2)
	SPARSEFLAG_INCOMPAT_FLAGS         uint32 = 0xFFFF0000
	SPARSEFLAG_COMPRESSED                    = (1 << 16)
	SPARSEFLAG_EMBEDDED_LBA                  = (1 << 17)
)

const VMDK_SECTOR_SIZE uint64 = 512

const (
	GRAIN_MARKER_EOS = iota
	GRAIN_MARKER_GRAIN_TABLE
	GRAIN_MARKER_GRAIN_DIRECTORY
	GRAIN_MARKER_FOOTER
	GRAIN_MARKER_PROGRESS
)

type (
	LEUint16 []byte
	LEUint32 []byte
	LEUint64 []byte
)

var (
	sizeOfUint16 = binary.Size(uint16(0))
	sizeOfUint32 = binary.Size(uint32(0))
	sizeOfUint64 = binary.Size(uint64(0))
)

func cpuUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func cpuUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func cpuUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func leUint16(u uint16) []byte {
	b := make([]byte, sizeOfUint16)
	binary.LittleEndian.PutUint16(b, u)
	return b
}

func leUint32(u uint32) []byte {
	b := make([]byte, sizeOfUint32)
	binary.LittleEndian.PutUint32(b, u)
	return b
}

func leUint64(u uint64) []byte {
	b := make([]byte, sizeOfUint64)
	binary.LittleEndian.PutUint64(b, u)
	return b
}

type SparseExtentHeaderOnDisk struct {
	MagicNumber        LEUint32
	Version            LEUint32
	Flags              LEUint32
	Capacity           LEUint64
	GrainSize          LEUint64
	DescriptorOffset   LEUint64
	DescriptorSize     LEUint64
	NumGTEsPerGT       LEUint32
	RgdOffset          LEUint64
	GdOffset           LEUint64
	OverHead           LEUint64
	UncleanShutdown    uint8
	SingleEndLineChar  byte
	NonEndLineChar     byte
	DoubleEndLineChar1 byte
	DoubleEndLineChar2 byte
	CompressAlgorithm  LEUint16
	pad                [433]uint8
}

type SparseGrainLBAHeaderOnDisk struct {
	LBA     LEUint64
	CmpSize LEUint32
}

type SparseSpecialLBAHeaderOnDisk struct {
	LBA     LEUint64
	CmpSize LEUint32
	Type    LEUint32
}

type SectorType uint64

type SparseExtentHeader struct {
	Version           uint32
	Flags             uint32
	NumGTEsPerGT      uint32
	CompressAlgorithm uint16
	UncleanShutdown   uint8
	Reserved          uint8
	Capacity          SectorType
	GrainSize         SectorType
	DescriptorOffset  SectorType
	DescriptorSize    SectorType
	RgdOffset         SectorType
	GdOffset          SectorType
	OverHead          SectorType
}

func bitwiseTrue(a, b uint32) bool {
	return a&b > 0
}

func ceiling(x, y uint64) uint64 {
	return (((x) + (y) - 1) / (y))
}

var mrand48 = rand.New(rand.NewSource(time.Now().UnixNano()))
