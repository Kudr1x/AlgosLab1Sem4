package algorithms

type RLE struct{}

func (r *RLE) Encode(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return []byte{}, nil
	}

	var result []byte
	i := 0

	for i < len(input) {
		current := input[i]
		j := i + 1

		for j < len(input) && input[j] == current && (j-i) < 127 {
			j++
		}

		if j-i > 1 {
			result = append(result, byte(j-i), current)
			i = j
			continue
		}

		start := i
		j = i + 1
		maxLen := min(len(input), start+127)

		for j < maxLen {
			if j+1 < maxLen && input[j] != input[j+1] {
				j++
			} else {
				break
			}
		}
		literalLen := j - start

		result = append(result, byte(0x80|literalLen))
		result = append(result, input[start:start+literalLen]...)
		i = j
	}

	return result, nil
}

func (r *RLE) Decode(encoded []byte) ([]byte, error) {
	var result []byte
	i := 0

	for i < len(encoded) {
		if i >= len(encoded) {
			break
		}
		control := encoded[i]
		i++

		if (control & 0x80) != 0 {
			length := int(control & 0x7F)
			end := i + length
			if end > len(encoded) {
				end = len(encoded)
			}
			result = append(result, encoded[i:end]...)
			i = end
		} else {
			if i >= len(encoded) {
				break
			}
			value := encoded[i]
			i++
			count := int(control)
			for j := 0; j < count; j++ {
				result = append(result, value)
			}
		}
	}

	return result, nil
}
