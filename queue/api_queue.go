package queue

import (
	"wa-blast/components"
	"wa-blast/components/rabbitmq"
	"wa-blast/loggers"
)

// Components
var mq *rabbitmq.RabbitComponent

// Logger
var log = loggers.Get()

// Init ...
func Init() {
	// Get component instance
	mq = components.GetRabbitMQ("mq")

	// Init queue routing
	mq.DeclareQueues()

	// Init subscription

}
