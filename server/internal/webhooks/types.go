package webhook

import "github.com/luponetn/nexusmq/pkg/broker"

type receiveResult struct {
	msg *broker.Message
	err error
}
