package readwrite

import (
	"encoding/binary"
	"io"
	"os"
)

type Writer struct {
	file *os.File
}

type FileEntry struct {
	FileName      string
	FileNameLower uint32
	FileNameUpper uint32
	Offset        uint64
	UncompSize    uint64
}

type DataEntry struct {
	Hash     uint32
	FileName string
}

func FindByHash(data []DataEntry, hash uint32) *DataEntry {
	for _, entry := range data {
		if entry.Hash == hash {
			return &entry
		}
	}
	return nil
}

func FindByFileName(data []DataEntry, fileName string) *DataEntry {
	for _, entry := range data {
		if entry.FileName == fileName {
			return &entry
		}
	}
	return nil
}

func NewWriter(fileName string, append bool) (*Writer, error) {
	var file *os.File
	var err error
	if append {
		file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	} else {
		file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	}
	if err != nil {
		return nil, err
	}
	return &Writer{file}, nil
}

func (w *Writer) WriteUInt32(value uint32) error {
	return binary.Write(w.file, binary.LittleEndian, value)
}

func (w *Writer) WriteUInt64(value uint64) error {
	return binary.Write(w.file, binary.LittleEndian, value)
}

func (w *Writer) Write(data []byte) (int, error) {
	return w.file.Write(data)
}

func (w *Writer) WriteChar(data string) (int, error) {
	return w.file.WriteString(data)
}

func (w *Writer) Seek(position int64, whence int) (int64, error) {
	return w.file.Seek(position, whence)
}

func (w *Writer) SeekFromBeginning(position int64) (int64, error) {
	return w.file.Seek(position, io.SeekStart)
}

func (w *Writer) SeekFromEnd(position int64) (int64, error) {
	return w.file.Seek(position, io.SeekEnd)
}

func (w *Writer) SeekFromCurrent(position int64) (int64, error) {
	return w.file.Seek(position, io.SeekCurrent)
}

func (w *Writer) Position() (int64, error) {
	return w.file.Seek(0, io.SeekCurrent)
}

func (w *Writer) Size() (int64, error) {
	currentPos, err := w.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	defer w.file.Seek(currentPos, io.SeekStart)
	fileSize, err := w.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	return fileSize, nil
}

func (w *Writer) Close() error {
	return w.file.Close()
}
