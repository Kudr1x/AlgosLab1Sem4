package algorithms

import (
	"bytes"
	"container/heap"
	"encoding/binary"
	"fmt"
	"io"
)

type HA struct{}

type huffmanNode struct {
	char        byte
	freq        int
	left, right *huffmanNode
}

type PriorityQueue []*huffmanNode

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].freq < pq[j].freq }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*huffmanNode))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func buildFrequencyTable(data []byte) [256]uint32 {
	var freqTable [256]uint32
	for _, b := range data {
		freqTable[b]++
	}
	return freqTable
}

func buildTree(freqTable [256]uint32) *huffmanNode {
	pq := &PriorityQueue{}
	heap.Init(pq)

	for i, f := range freqTable {
		if f > 0 {
			heap.Push(pq, &huffmanNode{char: byte(i), freq: int(f)})
		}
	}

	if pq.Len() == 0 {
		return &huffmanNode{}
	}

	for pq.Len() > 1 {
		left := heap.Pop(pq).(*huffmanNode)
		right := heap.Pop(pq).(*huffmanNode)

		parent := &huffmanNode{
			freq:  left.freq + right.freq,
			left:  left,
			right: right,
		}
		heap.Push(pq, parent)
	}
	return heap.Pop(pq).(*huffmanNode)
}

type BitWriter struct {
	buffer []byte
	offset int
}

func (bw *BitWriter) WriteBits(code string) {
	for _, c := range code {
		if bw.offset%8 == 0 {
			bw.buffer = append(bw.buffer, 0)
		}
		pos := bw.offset / 8
		bitPos := 7 - (bw.offset % 8)
		if c == '1' {
			bw.buffer[pos] |= 1 << bitPos
		}
		bw.offset++
	}
}

func generateCodes(root *huffmanNode) map[byte]string {
	codes := make(map[byte]string)
	var traverse func(*huffmanNode, string)
	traverse = func(node *huffmanNode, code string) {
		if node == nil {
			return
		}
		if node.left == nil && node.right == nil {
			codes[node.char] = code
			return
		}
		traverse(node.left, code+"0")
		traverse(node.right, code+"1")
	}
	traverse(root, "")
	return codes
}

func (h *HA) Encode(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}

	freqTable := buildFrequencyTable(data)
	root := buildTree(freqTable)
	codes := generateCodes(root)

	bw := &BitWriter{}
	for _, b := range data {
		code, ok := codes[b]
		if !ok {
			return nil, fmt.Errorf("нет кода для символа: %v", b)
		}
		bw.WriteBits(code)
	}

	var outBuf bytes.Buffer
	for i := 0; i < 256; i++ {
		binary.Write(&outBuf, binary.LittleEndian, freqTable[i])
	}

	padding := byte((8 - (bw.offset % 8)) % 8)
	outBuf.WriteByte(padding)
	outBuf.Write(bw.buffer)

	return outBuf.Bytes(), nil
}

type BitReader struct {
	data   []byte
	offset int
}

func (br *BitReader) ReadBit() (byte, error) {
	if br.offset >= len(br.data)*8 {
		return 0, io.EOF
	}
	pos := br.offset / 8
	bitPos := 7 - (br.offset % 8)
	br.offset++
	return (br.data[pos] >> bitPos) & 1, nil
}

func (h *HA) Decode(data []byte) ([]byte, error) {
	if len(data) < 256*4+1 {
		return nil, fmt.Errorf("invalid encoded data")
	}

	var freqTable [256]uint32
	buf := bytes.NewReader(data)
	for i := 0; i < 256; i++ {
		binary.Read(buf, binary.LittleEndian, &freqTable[i])
	}

	padding, _ := buf.ReadByte()
	encodedData, _ := io.ReadAll(buf)

	root := buildTree(freqTable)
	br := &BitReader{data: encodedData}

	var outBuf bytes.Buffer
	node := root
	totalBits := len(encodedData)*8 - int(padding)

	for i := 0; i < totalBits; i++ {
		bit, err := br.ReadBit()
		if err != nil {
			break
		}

		if bit == 0 {
			node = node.left
		} else {
			node = node.right
		}

		if node.left == nil && node.right == nil {
			outBuf.WriteByte(node.char)
			node = root
		}
	}
	return outBuf.Bytes(), nil
}
