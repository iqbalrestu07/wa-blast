package rabbitmq

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"

	"os"

	"wa-blast/loggers"
)

const (
	// RabbitMQ key in components.yml file
	ComponentType = "rabbitmq"
	// Exchange Type
	ExchangeTypeDelayedMessage = "x-delayed-message"
)

// Ack
const (
	Ack = iota
	AckAll
	Ignore
	IgnoreRequeue
	IgnoreAll
	IgnoreRequeueAll
	Reject
	Requeue
)

var log = loggers.Get()

type RabbitComponent struct {
	producerConn *amqp.Connection
	consumerConn *amqp.Connection
	// TODO Add generic logger fields
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// GetConsumer returns consumer connection to implement advanced Consumer implementation
func (c *RabbitComponent) GetConsumer() *amqp.Connection {
	return c.consumerConn
}

// GetProducer returns producer connection to implement advanced Producer implementation
func (c *RabbitComponent) GetProducer() *amqp.Connection {
	return c.producerConn
}

// Queue represent structure for queue configuration
type Queue struct {
	Name            string
	Durable         bool
	AutoDelete      bool
	Exclusive       bool
	NoWait          bool
	Args            amqp.Table
	Exchange        *Exchange
	ExchangeBinding *QueueBinding
}

// NewQueue create a new instance of general usage Queue. It is durable, non-exclusive and not bound to an exchange
//
// For more customized usage, just init a new Queue struct, set ExchangeKey and RoutingKey to bind
func NewQueue(name string) Queue {
	return Queue{
		Name:    name,
		Durable: true,
	}
}

// NewDelayedQueue create a new instance of general usage Queue that bound to a delayed exchange.
// It is durable, non-exclusive and bound to exchange with routing key
//
// For more customized usage, just init a new Queue struct
func NewDelayedQueue(name, exchangeName, routingKey string) Queue {
	return Queue{
		Name:            name,
		Durable:         true,
		Exchange:        NewDelayedExchange(exchangeName),
		ExchangeBinding: NewQueueBinding(routingKey),
	}
}

// QueueBinding represent structure of queue binding to exchange configuration
type QueueBinding struct {
	RoutingKey string
	NoWait     bool
	Args       amqp.Table
}

// NewQueueBinding create new instance of general usage Queue Binding that is no wait has not args
func NewQueueBinding(routingKey string) *QueueBinding {
	return &QueueBinding{RoutingKey: routingKey}
}

// Exchange represent structure of exchange configuration
type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// NewDelayedExchange create new configuration for a durable Exchange that supports delayed message sending
//
// For a more customized delayed exchange, just initiate new struct and use NewDelayedExchangeArgs() func helper
// to generate a delayed exchange args
func NewDelayedExchange(name string) *Exchange {
	return &Exchange{
		Name:    name,
		Kind:    ExchangeTypeDelayedMessage,
		Durable: true,
		Args:    NewDelayedExchangeArgs(),
	}
}

// NewDelayedExchangeArgs returns arguments for a delayed message exchange
func NewDelayedExchangeArgs() amqp.Table {
	// Init arguments
	exchangeArgs := make(amqp.Table)
	// Set delayed message exchange args
	exchangeArgs["x-delayed-type"] = "direct"
	// Return args
	return exchangeArgs
}

// DeclareQueues takes queues as input and declare queues in bulk
func (c *RabbitComponent) DeclareQueues(queues ...Queue) {
	// Validate queues
	if len(queues) < 1 {
		fmt.Println("Failed to init routing. queues must be set")
		os.Exit(28)
		return
	}
	// Open consumer channel
	conn := c.consumerConn
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Failed to establish channel to consumer. Error: %s\n", err)
		os.Exit(29)
		return
	}
	// Close channel on returns
	defer ch.Close()
	// Iterate queues
	for _, v := range queues {
		// Declare queue if not exist
		q, err := ch.QueueDeclare(v.Name, v.Durable, v.AutoDelete, v.Exclusive, v.NoWait, v.Args)
		if err != nil {
			if err != nil {
				fmt.Printf("Failed to declare queue. Error: %s\n", err)
				os.Exit(30)
			}
		}
		// If exchange is set, bind to exchange
		if v.Exchange != nil && v.ExchangeBinding != nil {
			x := v.Exchange
			r := v.ExchangeBinding
			// Declare exchange
			err := ch.ExchangeDeclare(x.Name, "direct", x.Durable, x.AutoDelete, x.Internal, x.NoWait, x.Args)
			if err != nil {
				fmt.Printf("Failed to declare exchange. Error: %s\n", err)
				os.Exit(31)
				return
			}
			// Bind queue to exchange
			err = ch.QueueBind(q.Name, r.RoutingKey, x.Name, r.NoWait, r.Args)
			if err != nil {
				fmt.Printf("Failed to bind queue to exchange. Error: %s\n", err)
				os.Exit(32)
				return
			}
		}
		log.Infof("Queue declared: %s", v.Name)
	}
}

// Publish message direct to queue
func (c *RabbitComponent) Publish(queueName string, payload interface{}) error {
	// Open channel
	conn := c.producerConn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// Close channel on return
	defer ch.Close()
	// Prepare payload
	bodyBytes, err := Encode(payload)
	if err != nil {
		return err
	}
	p := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         bodyBytes,
	}
	// Publish
	err = ch.Publish("", queueName, false, false, p)
	if err != nil {
		return err
	}
	log.Debugf("Message sent to queue: %s", queueName)
	return nil
}

// Publish message direct to queue
func (c *RabbitComponent) PublishRPC(queueName string, priority uint8, payload interface{}) (result string, err error) {
	rand.Seed(time.Now().UTC().UnixNano())

	// Open channel
	conn := c.producerConn
	ch, err := conn.Channel()
	if err != nil {
		return "", err
	}
	// Close channel on return
	defer ch.Close()
	// Prepare payload
	bodyBytes, err := Encode(payload)
	if err != nil {
		return "", err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	corrId := randomString(32)

	p := amqp.Publishing{
		ContentType:   "text/plain",
		CorrelationId: corrId,
		ReplyTo:       q.Name,
		Body:          bodyBytes,
		Priority:      priority,
	}

	// Publish
	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		p)

	if err != nil {
		return "", err
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			result = string(d.Body)
			break
		}
	}

	log.Debugf("Message sent to queue: %s", queueName)
	return
}

// PublishDelayed send messages to a delayed exchange with routing key. Set delay in millis
func (c *RabbitComponent) PublishDelayed(exchangeName string, routingKey string, payload interface{}, delay int64) error {
	// Open channel
	conn := c.producerConn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// Close channel on return
	defer ch.Close()
	// Set headers
	headers := make(amqp.Table)
	headers["x-delay"] = delay
	// Prepare payload
	bodyBytes, err := Encode(payload)
	if err != nil {
		return err
	}
	p := amqp.Publishing{
		Headers:      headers,
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         bodyBytes,
	}
	// Publish
	err = ch.Publish(exchangeName, routingKey, false, false, p)
	if err != nil {
		return err
	}
	log.Debugf("Message sent to delayed exchange: %s", exchangeName)
	return nil
}

// Subscribe to queue for general-purpose subscription. Queue to subscribe must not be exclusive, consumer
// will be unique and randomized and Ack is handled by returning Ack, NAck, Reject on ConsumerCallback function
//
// For more advanced subscription, use GetConsumer() to retrieve consumer connection, create channel manually and
// create custom handling
func (c *RabbitComponent) Subscribe(queueName string, fn ConsumerCallback) error {
	// Open channel
	conn := c.consumerConn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// Close channel on returns
	defer ch.Close()
	// Set QoS for channel
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Error(err)
		return err
	}
	// Start consuming queue
	messages, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	go func() {
		for d := range messages {
			log.Debugf("Received message for queue %s", queueName)
			if fn != nil {
				// Execute callback
				ack := fn(d.Body)
				switch ack {
				case AckAll:
					d.Ack(true)
				case Ignore:
					d.Nack(false, false)
				case IgnoreRequeue:
					d.Nack(false, true)
				case IgnoreAll:
					d.Nack(true, false)
				case IgnoreRequeueAll:
					d.Nack(true, true)
				case Reject:
					d.Reject(false)
				case Requeue:
					d.Reject(true)
				default:
					d.Ack(false)
				}
			} else {
				// Send ack for this message only
				d.Ack(false)
			}
		}
	}()
	// Reading from an empty channel to block function from returning
	// Source: https://blog.sgmansfield.com/2016/06/how-to-block-forever-in-go/
	log.Infof("%s receiver initiated", queueName)
	forever := make(chan bool)
	<-forever
	// Return success
	return nil
}

// Subscribe to queue for general-purpose subscription. Queue to subscribe must not be exclusive, consumer
// will be unique and randomized and Ack is handled by returning Ack, NAck, Reject on ConsumerCallback function
//
// For more advanced subscription, use GetConsumer() to retrieve consumer connection, create channel manually and
// create custom handling
func (c *RabbitComponent) SubscribeRPC(queueName string, fn ConsumerCallbackRPC) error {
	// Open channel
	conn := c.consumerConn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// Close channel on returns
	defer ch.Close()
	// Set QoS for channel
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Error(err)
		return err
	}
	// Start consuming queue
	messages, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Error(err)
		return err
	}

	go func() {
		for d := range messages {

			if fn != nil {
				// Execute callback
				data, _ := fn(d.Body)

				err = ch.Publish(
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(data),
					})
				if err != nil {
					log.Infof("ERROR Queue %s", queueName)
				}
				d.Ack(false)
			}

		}
	}()

	// Reading from an empty channel to block function from returning
	// Source: https://blog.sgmansfield.com/2016/06/how-to-block-forever-in-go/
	log.Infof("%s receiver initiated", queueName)
	forever := make(chan bool)
	<-forever
	// Return success
	return nil
}

// MustSubscribe exit app if app failed to subscribe to a topic
func (c *RabbitComponent) MustSubscribe(queueName string, fn ConsumerCallback) {
	err := c.Subscribe(queueName, fn)
	if err != nil {
		fmt.Errorf("unable to subscribe queue: %s", queueName)
		os.Exit(27)
	}
}

// MustSubscribe exit app if app failed to subscribe to a topic
func (c *RabbitComponent) MustSubscribeRPC(queueName string, fn ConsumerCallbackRPC) {
	err := c.SubscribeRPC(queueName, fn)
	if err != nil {
		fmt.Errorf("unable to subscribe queue: %s", queueName)
		os.Exit(27)
	}
}

func Init(username, password, host, port string) *RabbitComponent {
	// Generate url
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", username, password, host, port)
	// Init component
	var c RabbitComponent
	// Create publish connection
	c.producerConn = newConnection(url)
	log.Debug("RabbitMQ Producer connection established")
	// Create subscriber connection
	c.consumerConn = newConnection(url)
	log.Debug("RabbitMQ Consumer connection established")

	// Return
	return &c
}
