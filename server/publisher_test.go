package server

import (
	"testing"
)

func TestNewPublisherServer(t *testing.T) {
	testPubSer := NewPublisherServer()

	if err := testPubSer.CacheStore.Set("test", []byte("set")); err != nil {
		t.Log(err)
		t.Fail()
	}
}