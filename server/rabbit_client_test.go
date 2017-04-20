package server

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

const (
	RabbitDockerImage = "rabbitmq:3"
)

var (
	isRabbitReady      = false
	rabbitMaxHost      = 5
	rabbitSecretCookie = "darkSide"
	rabbitUser         = "bugsbunny"
	rabbitPassword     = "wascalywabbit"
	rabbitEnvVars      = []string{
		"RABBITMQ_ERLANG_COOKIE=" + rabbitSecretCookie,
		// go go gadget go faster.
		"RABBITMQ_HIPE_COMPILE=0",
		"RABBITMQ_DEFAULT_USER=" + rabbitUser,
		"RABBITMQ_DEFAULT_PASS=" + rabbitPassword,
	}
)

func TestMain(m *testing.M) {
	// we need to be able to cancel it so we can kill the rabbit container later
	// also, set up a deferred recovery so we can just kill the container if there
	// is a panic. basically, CLEAN UP YOUR SHIT.
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		if r := recover(); r != nil {
			cancel()
		}
	}()

	defer ctx.Done()

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	fmt.Printf("setting up rabbitmq docker container for test with these vars: %v\n", rabbitEnvVars)
	setupRabbitContainer(ctx, dockerClient)

	os.Exit(m.Run())
}

func setupRabbitContainer(ctx context.Context, cli *client.Client) {
	output, err := cli.ImagePull(ctx, RabbitDockerImage, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, output)

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: RabbitDockerImage,
			Env:   rabbitEnvVars,
		}, &container.HostConfig{
			NetworkMode: "host",
		}, nil, "rabbit-test-"+uuid.NewV4().String())
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func TestBasicRabbitSend(t *testing.T) {
	rabbitConn, err := NewRabbitClient(fmt.Sprintf("amqp://%s:%s@localhost:5672/", rabbitUser, rabbitPassword))
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

func TestBasicRabbitReceive(t *testing.T) {
	rabbitConn, err := NewRabbitClient(fmt.Sprintf("amqp://%s:%s@localhost:5672/", rabbitUser, rabbitPassword))
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

