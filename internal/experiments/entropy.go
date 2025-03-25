package experiments

import (
	"AlgosLab1Sem4/internal/compressors"
	"AlgosLab1Sem4/internal/util"
	"math"
	"os"
)

func calculateEntropy(data []byte) float64 {
	counts := make(map[byte]int)
	for _, b := range data {
		counts[b]++
	}

	var entropy float64
	length := float64(len(data))
	for _, count := range counts {
		probability := float64(count) / length
		entropy -= probability * math.Log2(probability)
	}

	return entropy
}

func processFile(filename string, blockSize int) (float64, error) {
	BWTMTFCompressor := compressors.NewBWTMTFCompressor(1024)

	data, _ := os.ReadFile(filename)

	totalEntropy := 0.0
	numBlocks := 0

	for i := 0; i < len(data); i += blockSize {
		end := i + blockSize
		if end > len(data) {
			end = len(data)
		}
		block := data[i:end]

		processedBlock, _ := BWTMTFCompressor.Compress(block)

		entropy := calculateEntropy(processedBlock)
		totalEntropy += entropy
		numBlocks++
	}

	if numBlocks == 0 {
		return 0.0, nil
	}

	return totalEntropy / float64(numBlocks), nil
}

func StartEntropy() {
	dataJson := "/home/kudrix/GolandProjects/AlgosLab1Sem4v2/results/entropy.json"
	blockSizes := []int{64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}

	results := make(map[int]float64)
	for _, bs := range blockSizes {
		avgEntropy, _ := processFile("/home/kudrix/GolandProjects/AlgosLab1Sem4v2/datasets/text/enwik7", bs)
		results[bs] = avgEntropy
	}

	util.SaveJson(dataJson, results)
	util.RunPython(dataJson)
}
