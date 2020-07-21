package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

var wg sync.WaitGroup
var TOPIC string = "Topic1"
var count int = 0

func consumer(id string) {
	fmt.Println("started reader")
	defer wg.Done()
	fmt.Println("Group-id ", id)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   TOPIC,
		GroupID: id,
	})
	fmt.Println("configration complete")
	fmt.Println("reader initialised")
	for {
		m, err := reader.ReadMessage(context.Background())
		fmt.Println(m)
		if err != nil {
			fmt.Println("some error occured", err)
			break
		}
		fmt.Println("Message is : ", string(m.Value))
		count++
	}
}

type response struct {
	Value int `json:"message_count"`
}

func messageCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := response{Value: count}
	//resp.Value = count
	// w.Write([]byte("OK"))
	js, _ := json.Marshal(resp)
	fmt.Println("MESSAGE COUNT:", count)
	fmt.Println("GET Response", resp)
	w.Write(js)
	//json.NewEncoder(w).Encode(resp)
	// //json.NewEncoder(w)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	groupid := strconv.Itoa(rand.Intn(999))
	router := mux.NewRouter()
	wg.Add(1)
	go consumer(groupid)
	router.HandleFunc("/mes_count", messageCount).Methods("GET")
	log.Fatal(http.ListenAndServe(":8001", router))
	wg.Wait()
}
