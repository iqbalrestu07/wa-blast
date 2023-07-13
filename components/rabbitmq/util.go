package rabbitmq

import (
	"fmt"
	"os"

	"bytes"
	"encoding/gob"

	"github.com/streadway/amqp"
)

// ConsumerCallback represents callback queue structure.
type ConsumerCallback func(p Payload) (ack int)
type ConsumerCallbackRPC func(p Payload) (data string, ack int)

// Payload extends byte array to be able to decode into struct
type Payload []byte

func (d Payload) Decode(dest interface{}) error {
	buf := bytes.NewBuffer(d)
	err := gob.NewDecoder(buf).Decode(dest)
	if err != nil {
		return err
	}
	return nil
}

// newConnection initiate connection to RabbitMQ
func newConnection(url string) *amqp.Connection {
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Printf("failed to connect to RabbitMQ. URL: %s, Error: %s\n", url, err.Error())
		os.Exit(26)
	}
	return conn
}

func Encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(dest interface{}, data []byte) error {
	buf := bytes.NewBuffer(data)
	err := gob.NewDecoder(buf).Decode(dest)
	if err != nil {
		return err
	}
	return nil
}
