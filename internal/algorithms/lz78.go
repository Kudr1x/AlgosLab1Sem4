package algorithms

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type LZ78 struct{}

func (l *LZ78) Encode(data []byte) ([]byte, error) {
	dict := make(map[string]uint32)
	dict[""] = 0
	nextIndex := uint32(1)

	var output bytes.Buffer
	var w []byte

	for _, b := range data {
		newW := append(w, b)
		if _, exists := dict[string(newW)]; exists {
			w = newW
		} else {
			prefixIndex := dict[string(w)]
			if err := binary.Write(&output, binary.LittleEndian, prefixIndex); err != nil {
				return nil, err
			}
			if err := output.WriteByte(1); err != nil {
				return nil, err
			}
			if err := output.WriteByte(b); err != nil {
				return nil, err
			}

			dict[string(newW)] = nextIndex
			nextIndex++
			w = []byte{}
		}
	}

	if len(w) > 0 {
		prefixIndex := dict[string(w)]
		if err := binary.Write(&output, binary.LittleEndian, prefixIndex); err != nil {
			return nil, err
		}
		if err := output.WriteByte(0); err != nil {
			return nil, err
		}
	}
	return output.Bytes(), nil
}

func (l *LZ78) Decode(data []byte) ([]byte, error) {
	dict := [][]byte{[]byte{}}
	reader := bytes.NewReader(data)
	var output bytes.Buffer

	for reader.Len() > 0 {
		var index uint32
		if err := binary.Read(reader, binary.LittleEndian, &index); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		flag, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		if flag == 1 {
			if reader.Len() < 1 {
				return nil, errors.New("неполный токен в закодированных данных")
			}
			sym, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			if int(index) >= len(dict) {
				return nil, errors.New("некорректный индекс в закодированных данных")
			}

			base := make([]byte, len(dict[index]))
			copy(base, dict[index])
			seq := append(base, sym)
			output.Write(seq)
			dict = append(dict, seq)
		} else if flag == 0 {
			if int(index) >= len(dict) {
				return nil, errors.New("некорректный индекс в закодированных данных")
			}
			seq := dict[index]
			output.Write(seq)
		} else {
			return nil, errors.New("недопустимое значение флага в закодированных данных")
		}
	}
	return output.Bytes(), nil
}
