package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

func serialize(msg interface{}) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}

func main() {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(
		logger.New(),
	)

	app.Get("/send", func(c *fiber.Ctx) error {
		data := make(map[string]string)
		data["param1"] = "Testing"
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        serialize(data),
		}

		if err := channelRabbitMQ.Publish(
			"",
			"QueueService1",
			false,
			false,
			message,
		); err != nil {
			return err
		}

		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
