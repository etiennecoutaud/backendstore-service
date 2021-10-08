package main

import (
	"encoding/json"
	"fmt"

	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueService struct {
	connectionString string
	ch               *amqp.Channel
	queueName        string
}

func newQueueService(name, password, endpoint string) *QueueService {
	return &QueueService{
		connectionString: fmt.Sprintf("amqp://%s:%s@%s", name, password, endpoint),
	}
}

func (qs *QueueService) onExit() {
	log.Println("close RabbitMQ connection")
	qs.ch.Close()
}

func (qs *QueueService) initConn() error {
	conn, err := amqp.Dial(qs.connectionString)
	if err != nil {
		log.Println("fail to establish connection with rabbitMq to %s, %s", qs.connectionString, err)
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println("fail to open channel", qs.connectionString, err)
		return err
	}
	qs.ch = ch
	log.Println("connection with RabbitMQ is established")
	return nil
}

func (qs *QueueService) declareNewQueue(queueName string) error {
	q, err := qs.ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Println("fail to create queue", qs.connectionString, err)
		return err
	}
	qs.queueName = q.Name
	log.Println("queue %s successfully created", queueName)
	return nil
}

func (qs *QueueService) sendMsg(item *Item) error {
	body, err := json.Marshal(item)
	if err != nil {
		log.Println("fail to unmarshal %v", item)
		return err
	}
	err = qs.ch.Publish(
		"",           // exchange
		qs.queueName, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Println("fail to send message")
		return err
	}
	log.Println("message sended")
	return nil
}
