package readwrite

import (
	"encoding/binary"
	"io"
	"os"
)

type Reader struct {
	file *os.File
}

func NewReader(fileName string) (*Reader, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return &Reader{file}, nil
}

func (r *Reader) IsValid() bool {
	return r.file != nil
}

func (r *Reader) ReadUInt32() (uint32, error) {
	var value uint32
	err := binary.Read(r.file, binary.LittleEndian, &value)
	return value, err
}

func (r *Reader) ReadUInt64() (uint64, error) {
	var value uint64
	err := binary.Read(r.file, binary.LittleEndian, &value)
	return value, err
}

func (r *Reader) Read(data []byte) (int, error) {
	return r.file.Read(data)
}

func (r *Reader) ReadChar() (byte, error) {
	var value byte
	err := binary.Read(r.file, binary.LittleEndian, &value)
	return value, err
}

func (r *Reader) Seek(position int64, whence int) (int64, error) {
	return r.file.Seek(position, whence)
}

func (r *Reader) SeekFromBeginning(position int64) (int64, error) {
	return r.file.Seek(position, io.SeekStart)
}

func (r *Reader) SeekFromEnd(position int64) (int64, error) {
	return r.file.Seek(position, io.SeekEnd)
}

func (r *Reader) SeekFromCurrent(position int64) (int64, error) {
	return r.file.Seek(position, io.SeekCurrent)
}

func (r *Reader) Position() (int64, error) {
	return r.file.Seek(0, io.SeekCurrent)
}

func (r *Reader) Size() (int64, error) {
	currentPos, err := r.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	defer r.file.Seek(currentPos, io.SeekStart)
	fileSize, err := r.file.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}
	return fileSize, nil
}

func (r *Reader) Close() error {
	return r.file.Close()
}
