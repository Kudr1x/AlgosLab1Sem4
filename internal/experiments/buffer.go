package experiments

import (
	"AlgosLab1Sem4/internal/compressors"
	"AlgosLab1Sem4/internal/util"
	"os"
)

func StartBuffer() {
	dataJson := "/home/kudrix/GolandProjects/AlgosLab1Sem4/results/buffer.json"
	data, _ := os.ReadFile("/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/color.png")

	bufferSizes := []int{64, 128, 256, 512, 1024, 2048, 4096}
	results := make(map[int]float64)

	for _, bufferSize := range bufferSizes {
		lz77 := compressors.NewLZ77Compressor(bufferSize)
		compressed, _ := lz77.Compress(data)
		compressionRatio := util.CalculateCompressionRatio(data, compressed)
		results[bufferSize] = compressionRatio
	}

	util.SaveJson(dataJson, results)
	util.RunPython(dataJson)
}
