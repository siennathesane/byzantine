package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/streadway/amqp"
)

const (
	RabbitVersion = 3
)

func TestDockerAvailability(t *testing.T) {
	httpClient := http.DefaultClient
	docker, err := client.NewClient(client.DefaultDockerHost, "1.13", httpClient, nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	//docker.ContainerCreate(context.Background(), types.ContainerCreateConfig{})
}

func TestRabbitSend(t *testing.T) {
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Log(err)
		t.Log("cannot connect to a local rabbit node!")
		t.FailNow()
	}
	defer rabbitConn.Close()

	rabbitChan, err := rabbitConn.Channel()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer rabbitChan.Close()

	rabbitQueue, err := rabbitChan.QueueDeclare(
		"testRabbitAvailability",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	testData := []byte("rabbitAvailabilityTest")
	err = rabbitChan.Publish(
		"",
		rabbitQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        testData,
		},
	)
}

func TestRabbitReceive(t *testing.T) {
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		t.Log(err)
		t.Log("cannot connect to a local rabbit node!")
		t.FailNow()
	}
	defer rabbitConn.Close()

	rabbitChan, err := rabbitConn.Channel()
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer rabbitChan.Close()

	rabbitQueue, err := rabbitChan.QueueDeclare(
		"testRabbitAvailability",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	msgs, err := rabbitChan.Consume(
		rabbitQueue.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		if string(msg.Body) == "rabbitAvailabilityTest" {
			return
		}
	}
	t.Fail()
}
