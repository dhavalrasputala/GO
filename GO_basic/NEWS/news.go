package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type article struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}
type newsResponse struct {
	Status   string    `json:"status"`
	Articles []article `json:"articles"`
}

func main() {
	apikey := "7533b12307c84c3dab76ff23a90ff31d"
	query := "golang"
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&apiKey=%s", query, apikey)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var newsData newsResponse
	err = json.NewDecoder(resp.Body).Decode(&newsData)
	if err != nil {
		panic(err)
	}
	for _, article := range newsData.Articles {
		fmt.Println(article.Title + "-" + article.URL)
	}
}
