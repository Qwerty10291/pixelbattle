package utils

import (
	"encoding/binary"
	"fmt"
	"os"
)

type number interface {
	uint8 | uint16 | uint32 | uint64
}

func Make2D[T any](n, m int) [][]T {
    matrix := make([][]T, n)
    rows := make([]T, n*m)
    for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
        endRow := startRow + m
        matrix[i] = rows[startRow:endRow:endRow]
    }
    return matrix
}

func ReadMatrix[T number](file *os.File) ([][]T, error) {
	var width uint32
	err := binary.Read(file, binary.LittleEndian, &width)
	if err != nil {
		return nil, err
	}
	var height uint32
	err = binary.Read(file, binary.LittleEndian, &height)
	if err != nil {
		return nil, err
	}
	data := make([]T, 0, width * height)
	err = binary.Read(file, binary.LittleEndian, &data)
	if err != nil {
		return nil, err
	}
	if len(data) != cap(data) {
		return nil, fmt.Errorf("invalid data size. Expected %d values, got %d", cap(data), len(data))
	}
	dataSize := len(data)
	expectedSize := int(width) * int(height)
	if dataSize != expectedSize {
		return nil, fmt.Errorf("invalid data size. Expected %d bytes, got %d bytes", expectedSize, dataSize)
	}
	matrix := make([][]T, height)
	for i := 0; i < int(height); i++ {
		matrix[i] = data[i*int(width) : (i+1)*int(width)]
	}
	return matrix, nil
}

func WriteMatrix[T number](data [][]T, file *os.File) error {
	height := uint32(len(data))
	if height == 0{
		return fmt.Errorf("matrix empty")
	}

	width := uint32(len(data[0]))
	err := binary.Write(file, binary.LittleEndian, width)
	if err != nil {
		return err
	}
	err = binary.Write(file, binary.LittleEndian, height)
	if err != nil {
		return err
	}

	for i := 0; i < len(data); i++{
		if err := binary.Write(file, binary.LittleEndian, data[i]); err != nil{
			return err
		}
	}
	return nil
}