package broker

import (
	nexusmq "github.com/luponetn/nexusmq/pkg/broker"
)

func NewBroker() nexusmq.Broker {
	broker := nexusmq.NewBroker()

	broker.CreateTopic("usage.ingested")
	broker.CreateTopic("invoice.created")

	return broker
}
