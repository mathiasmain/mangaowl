package API

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)


func requestAndJsonfyResponse(url string, data interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"+
		" AppleWebKit/537.36 (KHTML, like Gecko)"+
		" Chrome/129.0.0.0 Safari/537.36")
	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.Status[:3] == "403" || res.Status[:3] == "429" {
		log.Fatalln("You are making more request per second that you can do. Slow down, feela.")
	}
	fmt.Println("\nStatus da resposta do request:", res.Status)
	return json.NewDecoder(res.Body).Decode(data)
}