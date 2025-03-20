package main

import (
	"AlgosLab1Sem4/src/BWT"
	"AlgosLab1Sem4/src/RLE"
	"os"
)

func main() {
	inputFile := "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/txt/enwik5"
	encodedFile := "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/txt/encode"
	tempFile := "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/txt/temp"
	decodeFile := "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/txt/decode"

	BWT.EncodeWithBlocks(inputFile, encodedFile)
	encodeData, _ := os.ReadFile(encodedFile)
	encodeData = RLE.Encode(encodeData)
	os.WriteFile(encodedFile, encodeData, 0644)

	decodeData, _ := os.ReadFile(encodedFile)
	decodeData = RLE.Decode(decodeData)
	os.WriteFile(tempFile, decodeData, 0644)
	BWT.DecodeWithBlocks(tempFile, decodeFile)
}
