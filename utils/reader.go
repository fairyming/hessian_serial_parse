package utils

import (
	"errors"
	"fmt"
)

var ErrOutOfBounds = errors.New("operation is out of bounds")

type DataReader struct {
	data   []byte
	cursor int
}

func NewDataReader(data []byte) (*DataReader, error) {
	if data == nil {
		return nil, fmt.Errorf("DATA IS NIL")
	}
	return &DataReader{data: data, cursor: 0}, nil
}

// 检查剩余的数据是否足够 n 字节
func (d *DataReader) Enough(n int) bool {
	return d.cursor+n <= len(d.data)
}

// 当前是否已经读取到末尾
func (d *DataReader) IsEof() bool {
	return d.cursor == len(d.data)
}

// 返回当前数据指针的索引，返回值指向的字符还未读取过
func (d *DataReader) Position() int {
	return d.cursor
}

// peek
// 读取 n 个字节并返回，但是不移动指针
// 如果剩余的数据不够，则返回错误
func (d *DataReader) Peek(n int) ([]byte, error) {
	if d.Enough(n) {
		return d.data[d.cursor : d.cursor+n], nil
	} else {
		return nil, ErrOutOfBounds
	}
}

// 读取 n 个字节，如果不够则返回错误
func (d *DataReader) Read(n int) ([]byte, error) {
	if d.Enough(n) {
		result := d.data[d.cursor : d.cursor+n]
		d.cursor += n
		return result, nil
	} else {
		return nil, ErrOutOfBounds
	}
}

// 读取 n 个字节，如果不够则返回错误
func (d *DataReader) ReadString(n int) (string, error) {
	if d.Enough(n) {
		result := d.data[d.cursor : d.cursor+n]
		d.cursor += n
		return string(result), nil
	} else {
		return "", ErrOutOfBounds
	}
}

// 读取 n 个字节，如果不够则返回错误
func (d *DataReader) Read1() (byte, error) {
	if r, err := d.Read(1); err == nil {
		return r[0], nil
	} else {
		return 0, ErrOutOfBounds
	}
}

// 读取两个字节，并按照小端序转换为 u16
func (d *DataReader) ReadU16Le() (uint16, error) {
	tmp, err := d.Read(2)
	if err != nil {
		return 0, err
	}
	return uint16(tmp[0]) + uint16(tmp[1])<<8, err
}

// 读取两个字节，并按照大端序转换为 u16
func (d *DataReader) ReadU16Be() (uint16, error) {
	tmp, err := d.Read(2)
	if err != nil {
		return 0, err
	}
	return uint16(tmp[0])<<8 + uint16(tmp[1]), err
}

// 读取四个字节，并按照小端序转换为 u32
func (d *DataReader) ReadU32Le() (uint32, error) {
	tmp, err := d.Read(4)
	if err != nil {
		return 0, err
	}
	return uint32(tmp[0]) + uint32(tmp[1])<<8 + uint32(tmp[2])<<16 + uint32(tmp[3])<<24, err
}

// 读取四个字节，并按照大端序转换为 u32
func (d *DataReader) ReadU32Be() (uint32, error) {
	tmp, err := d.Read(4)
	if err != nil {
		return 0, err
	}
	return uint32(tmp[0])<<24 + uint32(tmp[1])<<16 + uint32(tmp[2])<<8 + uint32(tmp[3]), err
}

// 回退n字节
func (d *DataReader) Back(n int) error {
	if d.cursor < n {
		return fmt.Errorf("DATA BACK NOT ENOUGH")
	}
	d.cursor -= n
	return nil
}

// 跳过n字节
func (d *DataReader) Skip(n int) error {
	if d.Enough(n) {
		d.cursor += n
		return nil
	} else {
		return fmt.Errorf("DATA SKIP NOT ENOUGH")
	}
}
