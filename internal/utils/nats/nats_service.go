package nats

import (
	"authentication/internal/models"
	"encoding/json"
	"github.com/nats-io/nats.go"
)

type Service interface {
	RequestNotification(subject string, notification models.Notification) error
}

type natsService struct {
	nats string
}

func NewNatsService(nats string) Service {
	return &natsService{
		nats: nats,
	}
}

func (n *natsService) RequestNotification(subject string, notification models.Notification) error {
	conn, err := nats.Connect(n.nats)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	if err := conn.Publish(subject, data); err != nil {
		return err
	}

	return nil
}
