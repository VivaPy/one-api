package history

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/songquanpeng/one-api/relay/model"
	"time"
)

const queueName = "oneapi"

type MQRecorder struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
}

func NewMQRecorder() *MQRecorder {

	conn, err := amqp.Dial("amqp://golang:golang_rbmq1234@8.134.223.251:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "fail to create queue!")

	mqRecorder := &MQRecorder{
		conn:      conn,
		ch:        ch,
		queueName: q.Name,
	}
	return mqRecorder
}

func (mqr *MQRecorder) send(msg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSecond*time.Second)
	defer cancel()
	return mqr.ch.PublishWithContext(ctx,
		"",            // exchange
		mqr.queueName, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
}

func (mqr *MQRecorder) Push(tokenID int, messages []model.Message, usage *Usage) error {
	msgEntity := MessageToSend{tokenID, messages, *usage}
	msgBody, err := json.Marshal(&msgEntity)

	failOnError(err, "fail to marshal message")
	return mqr.send(msgBody)
}

func (mqr *MQRecorder) Pull(userID int) ([]model.Message, error) {
	panic("do not implemented!")
}
