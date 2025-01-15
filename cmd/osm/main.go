package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// URL для запроса остановок через Overpass API
	url := "http://overpass-api.de/api/interpreter?data=[out:json];(node[\"highway\"=\"bus_stop\"](42.95,47.45,42.99,47.55););out;"

	// Выполнение GET-запроса
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected status code: %d", resp.StatusCode)
	}

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Маршалинг
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Запись результата в файл JSON
	file, err := os.Create("stops.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(result)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}
}
