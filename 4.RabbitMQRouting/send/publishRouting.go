package main

import (
	"fmt"
	"rabbitmq_example/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	hjhOne := RabbitMQ.NewRabbitMQRouting("exOne", "key_one")
	hjhTwo := RabbitMQ.NewRabbitMQRouting("exOne", "key_two")

	for i:=0; i<10; i++ {
		hjhOne.PublishRouting("Hello hjh One!" + strconv.Itoa(i))
		hjhTwo.PublishRouting("Hello hjh Two!" + strconv.Itoa(i))
		time.Sleep(1*time.Second)
		fmt.Println(i)
	}
}