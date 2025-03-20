package RLE

import "bytes"

func EncodeWithCodeLength(input []byte, M int) []byte {
	symbolSize := M / 8
	if symbolSize == 0 || len(input)%symbolSize != 0 {
		return nil
	}

	var result []byte
	i := 0

	for i < len(input) {
		max := len(input) - symbolSize

		repeatCount := 1
		current := input[i : i+symbolSize]
		for j := i + symbolSize; j <= max; j += symbolSize {
			if bytes.Equal(current, input[j:j+symbolSize]) {
				if repeatCount == 127 {
					break
				}
				repeatCount++
			} else {
				break
			}
		}

		if repeatCount > 1 {
			result = append(result, byte(repeatCount))
			result = append(result, current...)
			i += repeatCount * symbolSize
			continue
		}

		nonRepeatCount := 0
		for k := i; k <= max; k += symbolSize {
			if k+symbolSize > max {
				nonRepeatCount++
				break
			}
			if bytes.Equal(input[k:k+symbolSize], input[k+symbolSize:k+2*symbolSize]) {
				break
			}
			nonRepeatCount++
			if nonRepeatCount == 127 {
				break
			}
		}

		if nonRepeatCount > 0 {
			result = append(result, byte(0x80|nonRepeatCount))
			result = append(result, input[i:i+nonRepeatCount*symbolSize]...)
			i += nonRepeatCount * symbolSize
		} else {
			result = append(result, 0x81)
			result = append(result, input[i:i+symbolSize]...)
			i += symbolSize
		}
	}

	return result
}

func DecodeWithCodeLength(input []byte, M int) []byte {
	symbolSize := M / 8
	if symbolSize == 0 {
		return nil
	}

	var result []byte
	i := 0

	for i < len(input) {
		control := input[i]
		i++
		count := int(control & 0x7F)
		isLiteral := (control & 0x80) != 0

		if isLiteral {
			end := i + count*symbolSize
			if end > len(input) {
				end = len(input)
			}
			result = append(result, input[i:end]...)
			i = end
		} else {
			if i+symbolSize > len(input) {
				break
			}
			symbol := input[i : i+symbolSize]
			i += symbolSize
			for j := 0; j < count; j++ {
				result = append(result, symbol...)
			}
		}
	}

	return result
}

func Encode(input []byte) []byte {
	if len(input) == 0 {
		return []byte{}
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

	return result
}

func Decode(encoded []byte) []byte {
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

	return result
}
