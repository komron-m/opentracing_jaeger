package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Publisher struct {
	conn     *amqp.Connection
	amqpChan *amqp.Channel
}

func NewPublisher() (*Publisher, error) {
	rabbitmqUser := os.Getenv("RABBITMQ_USER")
	rabbitmqPass := os.Getenv("RABBITMQ_PASSWORD")
	rabbitmqHost := os.Getenv("RABBITMQ_HOST")
	rabbitmqPort := os.Getenv("RABBITMQ_PORT")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitmqUser, rabbitmqPass, rabbitmqHost, rabbitmqPort)
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	return &Publisher{
		conn:     conn,
		amqpChan: ch,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, eventName string, payload any) error {
	// serialize payload/message to be sent
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// init message
	msg := amqp.Publishing{
		ContentType:     "text/json",
		ContentEncoding: "utf-8",
		DeliveryMode:    amqp.Persistent,
		Body:            body,
	}

	defaultExchange := os.Getenv("RABBITMQ_EXCHANGE_NAME")

	return p.amqpChan.Publish(defaultExchange, eventName, false, false, msg)
}

func (p *Publisher) PublishWithTraces(ctx context.Context, eventName string, payload any) error {
	// create a span
	span, ctx := opentracing.StartSpanFromContext(ctx, "Publisher.Publish", ext.SpanKindProducer)
	defer span.Finish()

	span.LogKV("event_name", eventName)

	// serialize span context
	bagItems := map[string]string{}
	if err := span.Tracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(bagItems)); err != nil {
		return err
	}
	bagItemsJsonBytes, err := json.Marshal(bagItems)
	if err != nil {
		return err
	}

	// serialize payload/message to be sent
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// init message
	msg := amqp.Publishing{
		Headers: amqp.Table{
			"opentracing_data": string(bagItemsJsonBytes),
		},
		ContentType:     "text/json",
		ContentEncoding: "utf-8",
		DeliveryMode:    amqp.Persistent,
		Body:            body,
	}

	return p.amqpChan.Publish(os.Getenv("RABBITMQ_EXCHANGE_NAME"), eventName, false, false, msg)
}
