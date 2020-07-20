package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

var wg sync.WaitGroup
var TOPIC string = "Topic1"
var chans [2]chan int
var count1 int

func produce(wg *sync.WaitGroup, message []byte) {
	defer wg.Done()
	fmt.Println("producer started ", TOPIC)
	partition := 0
	conn, _ := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", TOPIC, partition)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    TOPIC,
		Balancer: &kafka.LeastBytes{},
	})
	for {
		w.WriteMessages(context.Background(),
			kafka.Message{
				// Key:   []byte("Key-A"),
				Value: []byte(message),
			},
		)
	}
	conn.Close()
}

func consumer(wg *sync.WaitGroup, id string, ch chan int) {
	fmt.Println("started reader: ", id)
	defer wg.Done()
	conf := kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    TOPIC,
		GroupID:  id,
		MaxBytes: 10,
	}
	count := 0
	reader := kafka.NewReader(conf)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("some error occured", err)
			continue
		}

		fmt.Println("Message is : ", string(m.Value))
		count++
		ch <- count
	}
}

func sum() {
	fmt.Println("reading from Channel ")
	for {
		val1 := <-chans[1]
		val2 := <-chans[0]
		count1 = val1 + val2
	}
}
func triggerJson(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(req.Body)
	var result map[string]interface{}
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
func main() {

	router := mux.NewRouter()
	for id := 0; id < 2; id++ {
		chans[id] = make(chan int, 10)
		fmt.Println("starting consumer: ", id)
		wg.Add(1)
		go consumer(&wg, string(id), chans[id])
	}
	go sum()
	router.HandleFunc("/trigger", triggerJson).Methods("POST")
	router.HandleFunc("/home", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "NUMBER OF MESSAGES READ : %d", count1)
	})
	log.Fatal(http.ListenAndServe(":8000", router))

	wg.Wait()

}
