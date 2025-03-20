package BWT

type BlockHeader struct {
	Position uint64
	DataSize uint64
}

type rotation struct {
	Data []byte
	Num  int
}
