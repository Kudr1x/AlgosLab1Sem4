package experiments

import (
	"AlgosLab1Sem4/internal/compressors"
	"AlgosLab1Sem4/internal/util"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type CompressionResult struct {
	Compressor       string
	FileName         string
	OriginalSize     int
	CompressedSize   int
	DecompressedSize int
	CompressionRatio float64
}

func createCompressors() []compressors.Compressor {
	return []compressors.Compressor{
		compressors.NewHACompressor(),
		compressors.NewRLECompressor(),
		compressors.NewBWTRLECompressor(1024),
		compressors.NewBWTMTFHACompressor(1024),
		compressors.NewBWTMTFRLEHACompressor(1024),
		compressors.NewLZ77Compressor(4096),
		compressors.NewLZ77HACompressor(4096),
		compressors.NewLZ78Compressor(),
		compressors.NewLZ78HACompressor(),
	}
}

func createFiles() []string {
	return []string{
		"/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/text/enwik7",
		"/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/text/rutext",
		"/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/raw.bw",
		"/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/raw.color",
		"/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/image/raw.gray",
		"/home/kudrix/GolandProjects/AlgosLab1Sem4/datasets/bin/bin.apk",
	}
}

func StartCompression() {
	COMPRESSORS := createCompressors()
	FILES := createFiles()

	resultsCh := make(chan CompressionResult)
	errCh := make(chan error)
	var wg sync.WaitGroup

	maxGoroutines := runtime.NumCPU()
	sem := make(chan struct{}, maxGoroutines)

	for _, compressor := range COMPRESSORS {
		for _, file := range FILES {
			wg.Add(1)
			go func(c compressors.Compressor, f string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				res, err := processCompression(c, f)
				if err != nil {
					errCh <- err
					return
				}
				resultsCh <- res
			}(compressor, file)
		}
	}

	go func() {
		wg.Wait()
		close(resultsCh)
		close(errCh)
	}()

	var results []CompressionResult
	var errors []error

	for res := range resultsCh {
		results = append(results, res)
	}

	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		panic(errors[0])
	}

	writeCSV(results)
}

func processCompression(c compressors.Compressor, file string) (CompressionResult, error) {
	originalData, err := os.ReadFile(file)
	if err != nil {
		return CompressionResult{}, err
	}
	originalSize := len(originalData)

	compressedData, err := c.Compress(originalData)
	if err != nil {
		return CompressionResult{}, err
	}
	compressedSize := len(compressedData)

	compressedPath := strings.Replace(file[:strings.LastIndex(file, "/")+1]+c.GetName()+file[strings.LastIndex(file, "/")+1:], "datasets", "results/encode", 1)
	if err := os.MkdirAll(filepath.Dir(compressedPath), 0755); err != nil {
		return CompressionResult{}, err
	}
	if err := os.WriteFile(compressedPath, compressedData, 0644); err != nil {
		return CompressionResult{}, err
	}

	decompressedData, err := c.Decompress(compressedData)
	if err != nil {
		return CompressionResult{}, err
	}
	decompressedSize := len(decompressedData)

	decompressedPath := strings.Replace(file[:strings.LastIndex(file, "/")+1]+c.GetName()+file[strings.LastIndex(file, "/")+1:], "datasets", "results/decode", 1)
	if err := os.MkdirAll(filepath.Dir(decompressedPath), 0755); err != nil {
		return CompressionResult{}, err
	}
	if err := os.WriteFile(decompressedPath, decompressedData, 0644); err != nil {
		return CompressionResult{}, err
	}

	compressionRatio := util.CalculateCompressionRatio(originalData, compressedData)

	return CompressionResult{
		Compressor:       c.GetName(),
		FileName:         filepath.Base(file),
		OriginalSize:     originalSize,
		CompressedSize:   compressedSize,
		DecompressedSize: decompressedSize,
		CompressionRatio: compressionRatio,
	}, nil
}

func writeCSV(results []CompressionResult) {
	file, err := os.Create("/home/kudrix/GolandProjects/AlgosLab1Sem4/results/compression_stats.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{
		"Compressor",
		"File",
		"Original Size (bytes)",
		"Compressed Size (bytes)",
		"Decompressed Size (bytes)",
		"Compression Ratio (%)",
	}
	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	for _, res := range results {
		record := []string{
			res.Compressor,
			res.FileName,
			strconv.Itoa(res.OriginalSize),
			strconv.Itoa(res.CompressedSize),
			strconv.Itoa(res.DecompressedSize),
			fmt.Sprintf("%.3f", res.CompressionRatio),
		}
		if err := writer.Write(record); err != nil {
			panic(err)
		}
	}
}
