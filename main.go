package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
	"wb-project1/config"
	"wb-project1/database"
	"wb-project1/kafka"
	"wb-project1/processor"
)

func main() {
	db, err := database.InitPostgres()
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	err = database.InitSchema(db)
	if err != nil {
		log.Fatal("Error on creeating schema:", err)
	}

	cache, err := database.LoadAllOrders(db)
	if err != nil {
		log.Fatal("Failure on receiving cache:", err)
	}

	h := &processor.OrderHandler{
		DB:    db,
		Cache: cache,
	}
	go kafka.Consume(config.KafkaTopic, h.HandleKafkaMessage)

	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/order/")
		if data, ok := cache[id]; ok {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(data))
			return
		}

		row := db.QueryRow("SELECT data FROM orders WHERE order_uid = $1", id)
		var jsonData string
		err := row.Scan(&jsonData)
		if err != nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		cache[id] = jsonData
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jsonData))
	})

	fmt.Println("HTTP server running on :" + config.ServicePort)
	http.Handle("/", http.FileServer(http.Dir("static")))
	log.Fatal(http.ListenAndServe(":"+config.ServicePort, nil))
}
