package main

import (
	"net/http"
	"os"
	"time"

	"log"

	"github.com/google/uuid"
)

var p *Parameters

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Parameters struct {
	rabbitMQUser     string
	rabbitMQPasswd   string
	rabbitMQEndpoint string
}

type App struct {
	qs *QueueService
}

func (app *App) hello(w http.ResponseWriter, r *http.Request) {
	item := &Item{
		ID:    uuid.New().String(),
		Name:  "sweatshirt",
		Price: 10,
	}
	err := app.qs.sendMsg(item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write([]byte("send ok"))
	}
}

func main() {
	time.Sleep(20 * time.Second)
	qs := newQueueService(p.rabbitMQUser, p.rabbitMQPasswd, p.rabbitMQEndpoint)
	err := qs.initConn()
	if err != nil {
		os.Exit(1)
	}
	defer qs.onExit()
	err = qs.declareNewQueue("shipping")
	if err != nil {
		os.Exit(1)
	}

	app := &App{qs: qs}
	// Create a mux for routing incoming requests
	m := http.NewServeMux()
	// All URLs will be handled by this function
	m.HandleFunc("/", app.hello)

	// Create a server listening on port 8000
	s := &http.Server{
		Addr:    ":8000",
		Handler: m,
	}

	// Continue to process new requests until an error occurs
	log.Fatal(s.ListenAndServe())
}

func init() {
	p = &Parameters{}
	p.rabbitMQUser = os.Getenv("RABBITMQ_USER")
	p.rabbitMQPasswd = os.Getenv("RABBITMQ_PASSWD")
	p.rabbitMQEndpoint = os.Getenv("RABBITMQ_ENDPOINT")
}
