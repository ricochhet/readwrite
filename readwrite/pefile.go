/*
 * readwrite
 * Copyright (C) 2024 readwrite contributors
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package readwrite

import (
	"debug/pe"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"os"
)

var (
	errInvalidOffsetOrByteRange = errors.New("invalid offset or byte range")
	errSectionHeaderIsSizeZero  = errors.New("section header size is 0")
	errSectionIsNil             = errors.New("section is nil")
	errNoBytes                  = errors.New("no bytes")
)

// COFFHeader
// 0x50, 0x45 = PE
// COFF_START_BYTES_LEN == len(COFFStartBytes).
var COFFStartBytes = []byte{0x50, 0x45, 0x00, 0x00} //nolint:gochecknoglobals // wontfix

const (
	COFFStartBytesLen = 4
	COFFHeaderSize    = 20
)

// OptionalHeader64
// https://github.com/golang/go/blob/master/src/debug/pe/pe.go
// uint byte size of OptionalHeader64 without magic mumber(2 bytes) or data directory(128 bytes)
// OptionalHeader64 size is 240
// (110).
var OH64ByteSize = binary.Size(OptionalHeader64X110{}) //nolint:exhaustruct,gochecknoglobals // wontfix

// DataDirectory
// 16 entries * 8 bytes / entry.
const (
	DataDirSize      = 128
	DataDirEntrySize = 8
)

// SectionHeader32
// https://github.com/golang/go/blob/master/src/debug/pe/section.go
// uint byte size of SectionHeader32 without name(8 bytes) or characteristics(4 bytes)
// (28).
var SH32ByteSize = binary.Size(SectionHeader32X28{}) //nolint:exhaustruct,gochecknoglobals // wontfix

const (
	SH32EntrySize           = 64
	SH32NameSize            = 8
	SH32CharacteristicsSize = 4
)

// Data structure.
type Data struct {
	Bytes []byte
	PE    pe.File
}

// Section structure (.ooa).
type Section struct {
	ContentID   string
	OEP         uint64
	EncBlocks   []EncBlock
	ImageBase   uint64
	SizeOfImage uint32
	ImportDir   DataDir
	IATDir      DataDir
	RelocDir    DataDir
}

// Import structure (.ooa).
type Import struct {
	Characteristics uint32
	Timedatestamp   uint32
	ForwarderChain  uint32
	Name            uint32
	FThunk          uint32
}

// Thunk structure (.ooa).
type Thunk struct {
	Function uint32
	DataAddr uint32
}

// DataDir structure (.ooa).
type DataDir struct {
	VA   uint32
	Size uint32
}

// EncBlock structure (.ooa).
type EncBlock struct {
	VA          uint32
	RawSize     uint32
	VirtualSize uint32
	Unk         uint32
	CRC         uint32
	Unk2        uint32
	CRC2        uint32
	Pad         uint32
	FileOffset  uint32
	Pad2        uint64
	Pad3        uint32
}

type OptionalHeader64X110 struct {
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
}

type SectionHeader32X28 struct {
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLineNumbers uint32
	NumberOfRelocations  uint16
	NumberOfLineNumbers  uint16
}

func Open(path string) (*Data, error) {
	newData := new(Data)
	file, err := os.Open(path)
	if err != nil { //nolint:wsl // gofumpt conflict
		return nil, err
	}

	pefile, err := pe.NewFile(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	allBytes, err := io.ReadAll(file)
	if err != nil {
		file.Close()
		return nil, err
	}

	newData.Bytes = allBytes
	newData.PE = *pefile

	return newData, nil
}

func WriteBytes(bytes []byte, offset int, replace []byte) error {
	if offset < 0 || offset+len(replace) > len(bytes) {
		return errInvalidOffsetOrByteRange
	}

	copy(bytes[offset:], replace)

	return nil
}

func ReadCOFFHeaderOffset(bytes []byte) (int, error) {
	offset, err := FindBytes(bytes, COFFStartBytes)
	if err != nil {
		return -1, err
	}

	return offset, nil
}

func ReadDDBytes(bytes []byte) ([]byte, error) {
	offset, err := ReadCOFFHeaderOffset(bytes)
	if err != nil {
		return nil, err
	}

	return bytes[offset+COFFStartBytesLen+COFFHeaderSize+OH64ByteSize : offset+COFFStartBytesLen+COFFHeaderSize+OH64ByteSize+DataDirSize], nil
}

func ReadDDEntryOffset(bytes []byte, entryVirtualAddress uint32, entrySize uint32) (int, error) {
	dataDir, err := ReadDDBytes(bytes)
	if err != nil {
		return -1, err
	}

	entryBytes := make([]byte, DataDirEntrySize)
	binary.LittleEndian.PutUint32(entryBytes[:4], entryVirtualAddress)
	binary.LittleEndian.PutUint32(entryBytes[4:], entrySize)
	rva, err := FindBytes(dataDir, entryBytes)

	if err != nil || rva == -1 {
		log.Fatal(err)
		return -1, err
	}

	offset, err := ReadCOFFHeaderOffset(bytes)
	if err != nil {
		return -1, err
	}

	return offset + COFFStartBytesLen + COFFHeaderSize + OH64ByteSize + rva, nil
}

func ReadSHSize(file pe.File) (int, error) {
	sections := len(file.Sections)
	size := sections * SH32EntrySize

	if size == 0 {
		return -1, errSectionHeaderIsSizeZero
	}

	return size, nil
}

func ReadSHBytes(bytes []byte, shSize int) ([]byte, error) {
	offset, err := ReadCOFFHeaderOffset(bytes)
	if err != nil {
		return nil, err
	}

	return bytes[offset+COFFStartBytesLen+COFFHeaderSize+OH64ByteSize+DataDirSize : offset+COFFStartBytesLen+COFFHeaderSize+OH64ByteSize+DataDirSize+shSize], nil //nolint:lll // wontfix
}

func ReadSHEntryOffset(bytes []byte, address int) (int, error) {
	offset, err := ReadCOFFHeaderOffset(bytes)
	if err != nil {
		return -1, err
	}

	return offset + COFFStartBytesLen + COFFHeaderSize + OH64ByteSize + DataDirSize + address, nil
}

func ReadSectionBytes(file *Data, sectionVirtualAddress uint32, sectionSize uint32) ([]byte, error) {
	var section *pe.Section

	for _, s := range file.PE.Sections {
		if sectionVirtualAddress >= s.VirtualAddress && sectionVirtualAddress < s.VirtualAddress+s.Size {
			section = s
			break
		}
	}

	if section == nil {
		return nil, errSectionIsNil
	}

	offset := sectionVirtualAddress - section.VirtualAddress + section.Offset
	bytes := file.Bytes[offset : offset+sectionSize]

	return bytes, nil
}

func ReadImport(reader io.Reader) (Import, error) {
	var importData Import
	err := binary.Read(reader, binary.LittleEndian, &importData)

	return importData, err
}

func ReadThunk(reader io.Reader) (Thunk, error) {
	var thunkData Thunk
	err := binary.Read(reader, binary.LittleEndian, &thunkData)

	return thunkData, err
}

func ReadDataDir(reader io.Reader) (DataDir, error) {
	var dataDir DataDir
	err := binary.Read(reader, binary.LittleEndian, &dataDir)

	return dataDir, err
}

func ReadEncBlock(reader io.Reader) (EncBlock, error) {
	var encBlock EncBlock
	err := binary.Read(reader, binary.LittleEndian, &encBlock)

	return encBlock, err
}

func FindBytes(src []byte, dst []byte) (int, error) {
	for i := range src[:len(src)-len(dst)+1] {
		if MatchBytes(src[i:i+len(dst)], dst) {
			return i, nil
		}
	}

	return -1, errNoBytes
}

func PadBytes(bytes []byte, size int) []byte {
	if len(bytes) < size {
		paddingSize := size - len(bytes)
		padding := make([]byte, paddingSize)

		return append(bytes, padding...)
	}

	return bytes
}

func MatchBytes(src []byte, dst []byte) bool {
	for i := range dst {
		if src[i] != dst[i] {
			return false
		}
	}

	return true
}
