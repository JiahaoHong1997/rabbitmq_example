package main

import (
	"fmt"
	"rabbitmq_example/RabbitMQ"
	"time"
	"strconv"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("HongJiahaoWork")

	for i:=0; i<=100; i++ {
		rabbitmq.PublishSimple("Hello HongJiahao!" + strconv.Itoa(i))
		time.Sleep(1*time.Second)
		fmt.Println(i)
	}
}