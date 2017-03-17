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

func (ByzantineBrokerServer) RegisterSubscriber(context.Context, *byzantine.Subscriber) (*byzantine.ReadyResponse, error) {

	panic("implement me")
}

func (ByzantineBrokerServer) Ready(context.Context, *byzantine.Publication) (*byzantine.ReadyResponse, error) {
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
