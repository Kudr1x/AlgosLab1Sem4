package main

import (
	"AlgosLab1Sem4/src/IO"
	"AlgosLab1Sem4/src/RLE"
	"os"
)

func main() {
	inputFile := "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/txt/text.txt"
	encodedFile := "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/txt/txtencode"
	decodedFile := "/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/txt/txtdecode"

	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}

	M := 8
	encodedData := RLE.Encode(inputData, M)
	if encodedData == nil {
		panic("Ошибка кодирования")
	}

	if err := IO.WriteEncodedFile(encodedFile, encodedData, M); err != nil {
		panic(err)
	}

	readEncodedData, readM, err := IO.ReadEncodedFile(encodedFile)
	if err != nil {
		panic(err)
	}

	decodedData := RLE.Decode(readEncodedData, readM)
	if decodedData == nil {
		panic("Ошибка декодирования")
	}

	if err := os.WriteFile(decodedFile, decodedData, 0644); err != nil {
		panic(err)
	}
}
