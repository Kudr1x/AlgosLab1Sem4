package util

func CountingSort(S []byte) []int {
	N := len(S)
	M := 256
	T := make([]int, M)
	TSub := make([]int, M)

	for _, s := range S {
		T[s]++
	}

	for j := 1; j < M; j++ {
		TSub[j] = TSub[j-1] + T[j-1]
	}

	PInverse := make([]int, N)
	P := make([]int, N)

	for i := 0; i < N; i++ {
		PInverse[TSub[S[i]]] = i
		P[i] = TSub[S[i]]
		TSub[S[i]]++
	}

	return PInverse
}
