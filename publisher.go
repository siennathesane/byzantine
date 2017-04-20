package byzantine

import (
	"os"

	"github.com/satori/go.uuid"
)

func NewPublisher() (*Publisher, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return &Publisher{}, err
	}
	return &Publisher{
		Address:     hostname,
		PublisherID: uuid.NewV4().String(),
	}, nil
}

func (pub *Publisher) NewMessage(contents []byte) *Publication {
	return &Publication{
		Contents:      contents,
		Sender:        pub,
		PublicationID: uuid.NewV4().String(),
	}
}
