package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"sync"
)

type fileResult struct {
	Data  json.RawMessage `json:"data,omitempty"`
	Error string          `json:"error,omitempty"`
}

func uploadfile(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing multipart form", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["uploadfiles"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	results := make([]fileResult, len(files))
	var wg sync.WaitGroup

	for i, handler := range files {
		wg.Add(1)
		go func(h *multipart.FileHeader, idx int) {
			defer wg.Done()
			file, err := h.Open()
			if err != nil {
				results[idx] = fileResult{Error: err.Error()}
				return
			}

			csvReader := csv.NewReader(file)
			data, err := csvReader.ReadAll()
			if err != nil {
				results[idx] = fileResult{Error: err.Error()}
				return
			}
			file.Close()

			jsonData, err := CSVtoJSON(data)
			if err != nil {
				results[idx] = fileResult{Error: err.Error()}
				return
			}
			results[idx] = fileResult{Data: jsonData}
		}(handler, i)
	}
	wg.Wait()
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func CSVtoJSON(Data [][]string) ([]byte, error) {
	if len(Data) == 0 {
		return nil, fmt.Errorf("empty CSV file")
	}
	headers := Data[0]
	var records []map[string]string

	for _, line := range Data[1:] {
		rec := make(map[string]string)
		for j, field := range line {
			if j < len(headers) {
				rec[headers[j]] = field
			}
		}
		records = append(records, rec)
	}
	return json.Marshal(records)
}

func main() {
	http.HandleFunc("/upload", uploadfile)
	fmt.Println("Starting server at :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
