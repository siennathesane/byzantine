package server

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/mxplusb/byzantine"
	"google.golang.org/grpc/grpclog"
)

type ByzantinePublisherServer struct {
	*byzantine.Publisher
	CacheStore *bigcache.BigCache
	Logger     grpclog.Logger
}

func NewPublisherServer() *ByzantinePublisherServer {
	s := new(ByzantinePublisherServer)
	c, err := bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour * 6))
	if err != nil {
		grpclog.Fatalf("cannot instantiate publisher cache! %s", err)
	}
	s.CacheStore = c
	return s
}