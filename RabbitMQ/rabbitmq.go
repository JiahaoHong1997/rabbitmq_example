package RabbitMQ

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// url格式: amqp://账号:密码@rabbitmq服务器地址:端口号/vhost
const MQURL = "amqp://JiahaoHong1997:hjh1314lvwxy@127.0.0.1:5672/hjh"

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string // 队列名称
	Exchange  string // 交换机
	Key       string // key
	Mqurl     string // 连接信息
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}

	var err error
	// 创建rabbitmq连接
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "Connect Failed!")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "Fetch channel Failed!")
	return rabbitmq
}

// 断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 只使用默认交换机情况下的实例创建、发送和接收方法
// step1:创建Simple模式下的RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "") // Simple模式下只使用queque，是有默认交换机的（direct）
}

// step2:Simple模式下生产者代码
func (r *RabbitMQ) PublishSimple(message string) {
	// 1.申请队列，如果队列不存在会自动创建，如果存在则跳过创建。保证队列存在消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞
		nil,   // 额外属性
	)
	if err != nil {
		fmt.Println()
	}

	// 2.发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // 如果为true，根据exchange类型和routkey规则，如果无法找到符合条件的队列那么会吧发送的消息回退给发送者
		false, // 如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，则会吧消息发还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)

}

// step3:Simple模式下消费者代码
func (r *RabbitMQ) ConsumeSimple() {
	// 1.申请队列，如果队列不存在会自动创建，如果存在则跳过创建。保证队列存在消息能发送到队列中
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		false, // 是否持久化
		false, // 是否自动删除
		false, // 是否具有排他性
		false, // 是否阻塞
		nil,   // 额外属性
	)
	if err != nil {
		fmt.Println(err)
	}

	// 2.接受消息
	msgs, err := r.channel.Consume(
		r.QueueName,
		"",    // 用来区分多个消费者
		true,  // 是否自动应答
		false, // 是否具有排他性
		false, // 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false, // 队列消费是否阻塞（false为阻塞）
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启用协程处理消息
	go func() {
		for d := range msgs {
			// 实现我们要处理的逻辑函数
			log.Printf("Recieve a message: %s", d.Body)
		}
	}()

	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}

// 订阅模式
// 订阅模式下创建RabbitMQ实例
func NewRabbitMQPubSub(exchangeName string) *RabbitMQ {
	// 创建RabbitMQ实例
	rabbitmq := NewRabbitMQ("", exchangeName, "")
	return rabbitmq
}

// 订阅模式下生产
func (r *RabbitMQ) PublishPub(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout",
		true,
		false,
		false, // 如果设置为true，表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message)})
	r.failOnErr(err, "failed to publish message")
}

// 订阅模式消费端代码
func (r *RabbitMQ) RecieveSub() {
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"fanout", // 交换机类型
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare an exchange")

	// 2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生成队列名称
		false,
		false,
		true, // 排他性
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}


	// 3.绑定队列到交换机中
	err = r.channel.QueueBind(
		q.Name, // 队列名称
		"",     // 在pub/sub模式下，这里的key要为空
		r.Exchange,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}


	// 4.消费队列
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}


	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// 实现我们要处理的逻辑函数
			log.Printf("Recieve a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}


// Routing模式
// 创建RabbitMQ实例
func NewRabbitMQRouting(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("", exchangeName, routingKey)
	return rabbitmq
}

// Routing模式发送消息
func (r *RabbitMQ) PublishRouting(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err,"failed to drclare an exchange")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "test/plain",
			Body: []byte(message),})
	r.failOnErr(err, "failed to publish message!")
}

// 3.Routing模式接收消息
func (r *RabbitMQ) RecieveRouting() {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare a exchange!")

	// 2.尝试创建队列，这里的队列名称不要写
	q, err := r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare a queue!")

	// 3.绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 4.消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// 实现我们要处理的逻辑函数
			log.Printf("Recieve a message: %s", d.Body)
		}
	}()
	
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}

// Topic模式
// 创建RabbitMQ实例
func NewRabbitMQTopic(exchangeName string, routingKey string) *RabbitMQ {
	rabbitmq := NewRabbitMQ("",exchangeName,routingKey)
	return rabbitmq
}

// Topic模式发送消息
func (r *RabbitMQ) PublishTopic(message string) {
	// 1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "failed to declare an exchange!")

	// 2.发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),})

	r.failOnErr(err, "failed to publish message!")
}

// Topic模式接收消息
// 要注意key的规则
// 其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
// 匹配 hjh.* 表示匹配 hjh.hello,但是 hjh.hello.one 需要用 hjh.# 才能匹配到
func (r *RabbitMQ) RecieveTopic() {
	// 1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		r.Exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "failed to declare an exchange!")

	// 2.试探性创建队列
	q, err := r.channel.QueueDeclare(
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "failed to declare a queue!")

	// 3.绑定队列到exchange中
	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 4.消费消息
	msgs, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// 实现我们要处理的逻辑函数
			log.Printf("Recieve a message: %s", d.Body)
		}
	}()
	
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}