package algorithms

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type LZ77 struct {
	Buffer int
}

type lz77Token struct {
	Offset uint16
	Length uint16
	Next   byte
}

func (l *LZ77) Encode(data []byte) ([]byte, error) {
	var tokens []lz77Token
	n := len(data)
	i := 0

	for i < n {
		bestOffset, bestLength := 0, 0
		windowStart := 0
		if i > l.Buffer {
			windowStart = i - l.Buffer
		}

		for j := windowStart; j < i; j++ {
			length := 0
			for i+length < n &&
				j+length < i &&
				data[j+length] == data[i+length] &&
				length < 255 {
				length++
			}
			if length > bestLength {
				bestLength = length
				bestOffset = i - j
			}
		}

		var next byte
		if i+bestLength < n {
			next = data[i+bestLength]
		}

		tokens = append(tokens, lz77Token{
			Offset: uint16(bestOffset),
			Length: uint16(bestLength),
			Next:   next,
		})
		i += bestLength + 1
	}

	var outBuf bytes.Buffer
	for _, token := range tokens {
		binary.Write(&outBuf, binary.LittleEndian, token.Offset)
		binary.Write(&outBuf, binary.LittleEndian, token.Length)
		outBuf.WriteByte(token.Next)
	}
	return outBuf.Bytes(), nil
}

func (l *LZ77) Decode(data []byte) ([]byte, error) {
	tokenSize := 5
	if len(data)%tokenSize != 0 {
		return nil, fmt.Errorf("некорректные данные LZ77")
	}

	var tokens []lz77Token
	buf := bytes.NewReader(data)

	for buf.Len() > 0 {
		var token lz77Token
		binary.Read(buf, binary.LittleEndian, &token.Offset)
		binary.Read(buf, binary.LittleEndian, &token.Length)
		token.Next, _ = buf.ReadByte()
		tokens = append(tokens, token)
	}

	var result []byte
	for _, token := range tokens {
		start := len(result) - int(token.Offset)
		if start < 0 {
			return nil, fmt.Errorf("некорректное смещение")
		}

		for i := 0; i < int(token.Length); i++ {
			if start+i >= len(result) {
				return nil, fmt.Errorf("ошибка копирования данных")
			}
			result = append(result, result[start+i])
		}

		result = append(result, token.Next)
	}
	return result, nil
}
