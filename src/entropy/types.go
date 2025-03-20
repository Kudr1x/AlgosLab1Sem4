package entropy

type Data struct {
	Lengths   []int     `json:"lengths"`
	Entropies []float64 `json:"entropies"`
	Output    string    `json:"output"`
}
