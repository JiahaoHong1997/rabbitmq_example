package main

import (
	"fmt"
	"rabbitmq_example/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	hjhOne := RabbitMQ.NewRabbitMQTopic("exTopic", "hjh.topic.one")
	hjhTwo := RabbitMQ.NewRabbitMQTopic("exTopic", "hjh.topic.two")

	for i:=0; i<=10; i++ {
		hjhOne.PublishTopic("Hello hjh topic One!" + strconv.Itoa(i))
		hjhTwo.PublishTopic("Hello hjh topic Two!" + strconv.Itoa(i))
		time.Sleep(1*time.Second)
		fmt.Println(i)
	}
}