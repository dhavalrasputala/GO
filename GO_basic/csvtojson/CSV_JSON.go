package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
)

func uploadfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		w.Write([]byte("Error Handling your file "))
		return
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		w.Write([]byte("Error Handling your file "))
		return
	}

	jsondata, err := CSVtoJSON(data)
	if err != nil {
		w.Write([]byte("Error Handling your file "))
		return
	}
	w.Write([]byte(string(jsondata)))

	fmt.Fprintf(w, "File Converted successfully: %v", handler.Filename)
}
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
	http.HandleFunc("/upload", uploadfile)
	fmt.Println("Starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
