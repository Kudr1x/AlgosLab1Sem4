package algorithms

type MTF struct{}

func (m *MTF) Encode(data []byte) ([]byte, error) {
	alphabet := make([]byte, 256)
	for i := range alphabet {
		alphabet[i] = byte(i)
	}
	result := make([]byte, len(data))

	for i, b := range data {
		index := 0
		for ; index < len(alphabet) && alphabet[index] != b; index++ {
		}

		result[i] = byte(index)

		if index > 0 {
			tmp := alphabet[index]
			copy(alphabet[1:index+1], alphabet[0:index])
			alphabet[0] = tmp
		}
	}
	return result, nil
}

func (m *MTF) Decode(data []byte) ([]byte, error) {
	alphabet := make([]byte, 256)
	for i := range alphabet {
		alphabet[i] = byte(i)
	}
	result := make([]byte, len(data))

	for i, idx := range data {

		b := alphabet[idx]
		result[i] = b

		if idx > 0 {
			tmp := alphabet[idx]
			for j := idx; j > 0; j-- {
				alphabet[j] = alphabet[j-1]
			}
			alphabet[0] = tmp
		}
	}
	return result, nil
}
