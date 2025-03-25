package algorithms

type Algorithm interface {
	Encode(data []byte) ([]byte, error)
	Decode(data []byte) ([]byte, error)
}
