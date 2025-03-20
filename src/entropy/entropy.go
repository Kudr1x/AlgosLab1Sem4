package entropy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
)

func Calculate(inputFile, entropyVersusPlot, dataFile string) {
	data, _ := ioutil.ReadFile(inputFile)

	filteredData := filterASCII(data)

	lengths := []int{1, 2, 3, 4}
	entropies := make([]float64, len(lengths))

	for i, l := range lengths {
		entropies[i] = calculate(filteredData, l)
	}

	entropyData := Data{
		Lengths:   lengths,
		Entropies: entropies,
		Output:    entropyVersusPlot,
	}

	jsonData, _ := json.Marshal(entropyData)

	if err := ioutil.WriteFile(dataFile, jsonData, 0644); err != nil {
		fmt.Printf("Error writing data file: %v\n", err)
		os.Exit(1)
	}

	scriptPath := filepath.Join(".", "src/plot/plotter.py")
	cmd := exec.Command("python3", scriptPath, dataFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running Python script: %v\n", err)
		os.Exit(1)
	}
}

func calculate(data []byte, codeLength int) float64 {
	if codeLength < 1 || codeLength > 4 {
		return 0
	}

	adjustedLength := len(data) / codeLength * codeLength
	if adjustedLength == 0 {
		return 0
	}
	data = data[:adjustedLength]

	freq := make(map[[4]byte]int)
	total := 0

	for i := 0; i < adjustedLength; i += codeLength {
		var key [4]byte
		copy(key[:], data[i:i+codeLength])
		freq[key]++
		total++
	}

	if total == 0 {
		return 0
	}

	entropy := 0.0
	for _, count := range freq {
		p := float64(count) / float64(total)
		entropy += p * math.Log2(p)
	}

	return -entropy
}

func filterASCII(data []byte) []byte {
	var result []byte
	for _, b := range data {
		if b <= 127 {
			result = append(result, b)
		}
	}
	return result
}
