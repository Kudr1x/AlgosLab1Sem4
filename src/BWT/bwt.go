package BWT

import (
	"AlgosLab1Sem4/src/util"
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"sort"
)

const (
	blockSize = 1024 * 200
)

type BlockHeader struct {
	Position uint64
	DataSize uint64
}

type rotation struct {
	Data []byte
	Num  int
}

func Encode(S []byte) ([]byte, int) {
	N := len(S)
	BWM := make([][]byte, N)

	for i := 0; i < N; i++ {
		BWM[i] = append(S[i:], S[:i]...)
	}

	sort.Slice(BWM, func(i, j int) bool {
		return string(BWM[i]) < string(BWM[j])
	})

	lastColumn := make([]byte, N)
	for i := 0; i < N; i++ {
		lastColumn[i] = BWM[i][N-1]
	}

	var SIndex int
	for i := 0; i < N; i++ {
		if string(BWM[i]) == string(S) {
			SIndex = i
			break
		}
	}

	return lastColumn, SIndex
}

func Decode(lastColumn []byte, SIndex int) []byte {
	N := len(lastColumn)
	BWM := make([][]byte, N)

	for i := 0; i < N; i++ {
		BWM[i] = make([]byte, 0)
	}

	for step := 0; step < N; step++ {
		for j := 0; j < N; j++ {
			BWM[j] = append([]byte{lastColumn[j]}, BWM[j]...)
		}
		sort.Slice(BWM, func(i, j int) bool {
			return string(BWM[i]) < string(BWM[j])
		})
	}

	return BWM[SIndex]
}

func DecodeOptimized(lastColumn []byte, SIndex int) []byte {
	N := len(lastColumn)
	PInverse := util.CountingSort(lastColumn)

	S := make([]byte, 0)
	j := SIndex
	for i := 0; i < N; i++ {
		j = PInverse[j]
		S = append(S, lastColumn[j])
	}

	return S
}

func DecodeWithBlocks(inputPath, outputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	var blockCount uint64
	if err := binary.Read(input, binary.BigEndian, &blockCount); err != nil {
		return err
	}

	for i := uint64(0); i < blockCount; i++ {
		var header BlockHeader
		if err := binary.Read(input, binary.BigEndian, &header); err != nil {
			return err
		}

		data := make([]byte, header.DataSize)
		if _, err := io.ReadFull(input, data); err != nil {
			return err
		}

		restored := decodeBlock(data, int(header.Position))
		if _, err := output.Write(restored); err != nil {
			return err
		}
	}

	return nil
}

func EncodeWithBlocks(inputPath, outputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer output.Close()

	if err := binary.Write(output, binary.BigEndian, uint64(0)); err != nil {
		return err
	}

	var blockCount uint64
	buffer := make([]byte, blockSize)

	for {
		n, err := io.ReadFull(input, buffer)
		if err != nil && err != io.ErrUnexpectedEOF {
			if err == io.EOF {
				break
			}
			return err
		}

		transformed, pos := encodeBlock(buffer[:n])
		header := BlockHeader{
			Position: uint64(pos),
			DataSize: uint64(n),
		}

		if err := binary.Write(output, binary.BigEndian, header); err != nil {
			return err
		}

		if _, err := output.Write(transformed); err != nil {
			return err
		}

		blockCount++
	}

	if _, err := output.Seek(0, 0); err != nil {
		return err
	}

	return binary.Write(output, binary.BigEndian, blockCount)
}

func encodeBlock(data []byte) ([]byte, int) {
	if len(data) == 0 {
		return []byte{}, 0
	}

	rotations := make([]rotation, len(data))
	doubleData := append(data, data...)

	for i := range data {
		rotations[i] = rotation{
			Data: doubleData[i : i+len(data)],
			Num:  i,
		}
	}

	sort.Slice(rotations, func(i, j int) bool {
		return bytes.Compare(rotations[i].Data, rotations[j].Data) < 0
	})

	transformed := make([]byte, len(data))
	var position int

	for i, rot := range rotations {
		transformed[i] = rot.Data[len(data)-1]
		if rot.Num == 0 {
			position = i
		}
	}

	return transformed, position
}

func decodeBlock(data []byte, position int) []byte {
	if len(data) == 0 {
		return []byte{}
	}

	table := make([]int, 256)
	counts := make([]int, len(data))

	for i, b := range data {
		counts[i] = table[b]
		table[b]++
	}

	sum := 0
	prefix := make([]int, 256)
	for i := range prefix {
		prefix[i] = sum
		sum += table[i]
	}

	result := make([]byte, len(data))
	current := position
	result[len(result)-1] = data[current]

	for i := len(data) - 2; i >= 0; i-- {
		b := data[current]
		result[i] = b
		current = counts[current] + prefix[b]
	}

	return result
}
