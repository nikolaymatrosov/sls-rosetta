package main

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type ammo struct {
	Host    string            `json:"host"`
	Method  string            `json:"method"`
	Uri     string            `json:"uri"`
	Tag     string            `json:"tag"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func main() {
	// open file ammo.json
	file, err := os.Create("/Users/nikthespirit/GolandProjects/sls-rosetta/examples/go/concurrency/ammo.json")
	if err != nil {
		panic(err)
	}
	// defer close file
	defer file.Close()

	// in the loop from 0 to 10k generate json with random values and write it to file one per line
	for i := 0; i < 10000; i++ {
		j, err := json.Marshal(map[string]interface{}{
			"name": uuid.New().String(),
		})
		ammo := ammo{
			Host:   "functions.yandexcloud.net",
			Method: "POST",
			Uri:    "/d4e05r8jmt53mhnnurnb",
			Tag:    "long",
			Headers: map[string]string{
				"User-Agent":   "Load-Testing-Tool",
				"Connection":   "keep-alive",
				"X-Request-ID": uuid.New().String(),
				"Content-Type": "application/json",
			},
			Body: string(j),
		}

		jsonData, err := json.Marshal(ammo)
		if err != nil {
			panic(err)
		}

		_, err = file.WriteString(string(jsonData) + "\n")
		if err != nil {
			return
		}
	}
}
