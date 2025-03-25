package util

import (
	"encoding/json"
	"log"
	"os"
)

func SaveJson(fileName string, results map[int]float64) {
	jsonData, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
