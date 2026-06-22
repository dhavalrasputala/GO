package main

import (
	"encoding/json"
	"fmt"
)

func CSVtoJSON(Data [][]string) ([]byte, error) {
	headers := Data[0]
	var records []map[string]string

	for _, line := range Data[1:] {
		rec := make(map[string]string)
		for j, field := range line {
			rec[headers[j]] = field
		}
		records = append(records, rec)
	}
	return json.Marshal(records)
}

func main() {
	data := [][]string{
		{"Vegetable", "Fruit", "Rank"},
		{"Carrot", "Apple", "1"},
		{"Potato", "Banana", "2"},
	}

	jsonData, err := CSVtoJSON(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonData))
}
