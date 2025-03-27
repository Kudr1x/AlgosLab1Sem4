package algorithms

import (
	"encoding/binary"
	"errors"
	"sort"
)

type BWT struct {
	BlockSize int
}

func buildSuffixArray(s []byte) []int {
	n := len(s)
	sa := make([]int, n)
	rank := make([]int, n)
	temp := make([]int, n)

	for i := 0; i < n; i++ {
		sa[i] = i
		rank[i] = int(s[i])
	}

	k := 1
	cmp := func(i, j int) bool {
		if rank[i] != rank[j] {
			return rank[i] < rank[j]
		}
		ri, rj := -1, -1
		if i+k < n {
			ri = rank[i+k]
		}
		if j+k < n {
			rj = rank[j+k]
		}
		return ri < rj
	}

	for k < n {
		sort.Slice(sa, func(i, j int) bool {
			return cmp(sa[i], sa[j])
		})
		temp[sa[0]] = 0
		for i := 1; i < n; i++ {
			if cmp(sa[i-1], sa[i]) {
				temp[sa[i]] = temp[sa[i-1]] + 1
			} else {
				temp[sa[i]] = temp[sa[i-1]]
			}
		}
		copy(rank, temp)
		if rank[sa[n-1]] == n-1 {
			break
		}
		k *= 2
	}
	return sa
}

func bwtTransform(data []byte) ([]byte, uint32, error) {
	n := len(data)
	if n == 0 {
		return nil, 0, errors.New("empty block")
	}

	s := make([]byte, 2*n)
	copy(s, data)
	copy(s[n:], data)

	sa := buildSuffixArray(s)

	var filtered []int
	for _, idx := range sa {
		if idx < n {
			filtered = append(filtered, idx)
		}
	}

	var primary uint32
	for i, idx := range filtered {
		if idx == 0 {
			primary = uint32(i)
			break
		}
	}

	result := make([]byte, n)
	for i, idx := range filtered {
		if idx == 0 {
			result[i] = data[n-1]
		} else {
			result[i] = data[idx-1]
		}
	}

	return result, primary, nil
}

func inverseBWT(last []byte, primary uint32) ([]byte, error) {
	n := len(last)
	if n == 0 {
		return nil, errors.New("empty block")
	}
	count := make([]int, 256)
	for i := 0; i < n; i++ {
		count[last[i]]++
	}

	cumul := make([]int, 256)
	sum := 0
	for i := 0; i < 256; i++ {
		cumul[i] = sum
		sum += count[i]
	}

	next := make([]int, n)
	occ := make([]int, 256)
	for i := 0; i < n; i++ {
		symbol := last[i]
		next[i] = cumul[symbol] + occ[symbol]
		occ[symbol]++
	}

	res := make([]byte, n)
	idx := int(primary)
	for i := n - 1; i >= 0; i-- {
		res[i] = last[idx]
		idx = next[idx]
	}
	return res, nil
}

func (b *BWT) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}
	blockSize := b.BlockSize
	if blockSize <= 0 || blockSize > len(data) {
		blockSize = len(data)
	}
	var encoded []byte
	for i := 0; i < len(data); i += blockSize {
		end := i + blockSize
		if end > len(data) {
			end = len(data)
		}
		block := data[i:end]
		transformed, primary, err := bwtTransform(block)
		if err != nil {
			return nil, err
		}
		header := make([]byte, 8)
		binary.BigEndian.PutUint32(header[0:4], uint32(len(block)))
		binary.BigEndian.PutUint32(header[4:8], primary)
		encoded = append(encoded, header...)
		encoded = append(encoded, transformed...)
	}
	return encoded, nil
}

func (b *BWT) Decode(data []byte) ([]byte, error) {
	if len(data) < 8 {
		return nil, errors.New("input data is too short")
	}
	var decoded []byte
	pos := 0
	for pos < len(data) {
		if pos+8 > len(data) {
			return nil, errors.New("invalid header in encoded data")
		}
		blockLen := int(binary.BigEndian.Uint32(data[pos : pos+4]))
		primary := binary.BigEndian.Uint32(data[pos+4 : pos+8])
		pos += 8
		if pos+blockLen > len(data) {
			return nil, errors.New("block length exceeds data length")
		}
		block := data[pos : pos+blockLen]
		pos += blockLen
		original, err := inverseBWT(block, primary)
		if err != nil {
			return nil, err
		}
		decoded = append(decoded, original...)
	}
	return decoded, nil
}
