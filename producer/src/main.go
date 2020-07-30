package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

var wg sync.WaitGroup
var TOPIC string = "Topic1"

func produce(wg *sync.WaitGroup, message []byte) {
	defer wg.Done()
	fmt.Println("Producer started ", TOPIC)
	fmt.Println("connected")
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    TOPIC,
		Balancer: &kafka.LeastBytes{},
	})
	fmt.Println("configration done")
	for {
		fmt.Println(message)
		w.WriteMessages(context.Background(),
			kafka.Message{
				Value: []byte(message),
			},
		)
	}
}

func triggerJson(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	var result map[string]interface{}
	fmt.Println("waiting for json")
	err := decoder.Decode(&result)
	if err != nil {
		panic(err)
	}
	b, err := json.Marshal(result)
	fmt.Println(result)
	wg.Add(1)
	go produce(&wg, b)
	return
}

type response struct {
	Value int `json:"message_count"`
}

func totalCount(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	urls := []string{"http://consumer1:8001/mes_count", "http://consumer2:8001/mes_count"}
	sum := 0
	for _, url := range urls {
		fmt.Println(url)
		val, e := http.Get(url)
		body, _ := ioutil.ReadAll(val.Body)
		body_str := string(body)
		fmt.Println("Response:", body_str)
		if e != nil {
			fmt.Println("error occured", e)
		}
		res := response{}
		json.Unmarshal(body, &res)
		fmt.Println("VALUE", res.Value)
		sum = sum + res.Value
	}
	json.NewEncoder(w).Encode(response{sum})
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/trigger", triggerJson).Methods("POST")
	router.HandleFunc("/total_count", totalCount).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
	wg.Wait()

}
