package IO

import (
	"encoding/binary"
	"io"
	"os"
)

func WriteEncodedFile(filename string, data []byte, M int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := binary.Write(file, binary.LittleEndian, int32(M)); err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func ReadEncodedFile(filename string) ([]byte, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var M int32
	if err := binary.Read(file, binary.LittleEndian, &M); err != nil {
		return nil, 0, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	return data, int(M), nil
}
