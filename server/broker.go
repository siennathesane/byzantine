package server

import (
	"context"
	"errors"

	"github.com/mxplusb/byzantine"
	"google.golang.org/grpc/grpclog"
)

type ByzantineBrokerServer struct {
	subscribers map[uint32]*byzantine.Subscriber
	publishers  map[uint32]*byzantine.Publisher
	logger      grpclog.Logger
}

// Echo receives a blank message from the publisher as a heartbeat.
func (bz *ByzantineBrokerServer) Echo(ctx context.Context, pub *byzantine.Publication) (*byzantine.EchoResponse, error) {
	bz.logger.Printf("echo from %d", pub.PublisherID)
	return &byzantine.EchoResponse{
		Hello: true,
	}, nil
}

// GetSubscribers returns a list of the subscribers. The request must come from a subscribed member.
func (bz *ByzantineBrokerServer) GetSubscribers(sub *byzantine.Subscriber, stream byzantine.Broker_GetSubscribersServer) error {
	if _, ok := bz.subscribers[sub.PoolID]; !ok {
		bz.logger.Printf("%s not in pool", sub.Address)
		return errors.New("non-member asked for subscriber list")
	} else {
		for k := range bz.subscribers {
			stream.Send(bz.subscribers[k])
		}
		return nil
	}
}

// RegisterSubscriber registers a new subscriber to the broker server.
func (bz *ByzantineBrokerServer) RegisterSubscriber(ctx context.Context, sub *byzantine.Subscriber) (*byzantine.ReadyResponse, error) {
	if _, ok := bz.subscribers[sub.PoolID]; !ok {
		bz.logger.Printf("adding %s to the subscriber pool", sub.Address)
		bz.subscribers[sub.PoolID] = sub
		return &byzantine.ReadyResponse{Ready: true}, nil
	} else {
		bz.logger.Printf("%s:%d cannot be added to the subscriber pool because it already exists", sub.Address, sub.PoolID)
		return &byzantine.ReadyResponse{Ready: false}, errors.New("cannot add subscriber")
	}
}

func (bz *ByzantineBrokerServer) Ready(context.Context, *byzantine.Publication) (*byzantine.ReadyResponse, error) {
	/* TODO implement some ready logic. idk what that will look like yet.
		this will most likely ping other broker servers to see if they are ready.
		1. TODO: Validate the Subscriber thread is ready
		2. TODO: Get positive ready responses from the other brokers.
	*/
	panic("implement me")
}

func (ByzantineBrokerServer) Receive(context.Context, *byzantine.Publication) (*byzantine.PubResponse, error) {
	panic("implement me")
}

func (ByzantineBrokerServer) Push(byzantine.Broker_PushServer) error {
	panic("implement me")
}

func (ByzantineBrokerServer) Chain(context.Context, *byzantine.Publication) (*byzantine.ChainResponse, error) {
	panic("implement me")
}

type ByzantineSubscriberServer struct {
	*byzantine.Subscriber
	Logger *grpclog.Logger
}

// RegisterNewSubscriber will register a new subscriber with the Broker Server.
//func (bz *ByzantineBrokerServer) RegisterNewSubscriber(sub *byzantine.Subscriber) {
//	_, ok := bz.Subscribers[sub.PoolID]
//	if !ok {
//		bz.Logger.Printf("%s is already in the pool", sub.PoolID)
//	} else {
//		bz.Subscribers[sub.PoolID] = sub
//	}
//}
