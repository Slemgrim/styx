package queue

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Connection tba
type Connection struct {
	Connection *amqp.Connection
}

// Channel tba
type Channel struct {
	Instance *amqp.Channel
}

// Queue tba
type Queue struct {
	Instance *amqp.Queue
}

// NewConnection tba
func NewConnection(host string, port int, username string, password string) (*Connection, error) {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	log.Println(connectionString)
	connection, err := amqp.Dial(connectionString)

	if err != nil {
		return nil, err
	}

	return &Connection{connection}, nil
}

// Close closes the rabbit mq connection
func (q *Connection) Close() {
	q.Connection.Close()
}

// Channel creates a rabbit mq channel
func (q *Connection) Channel() (*Channel, error) {
	instance, err := q.Connection.Channel()

	if err != nil {
		return nil, err
	}

	return &Channel{instance}, nil
}

// DeclareQueue Declares a new queue or uses an existing queue with the given name and flags
func (c *Channel) DeclareQueue(name string, durable bool, autodelete bool, exclusive bool, noWait bool) (*Queue, error) {
	queue, err := c.Instance.QueueDeclare(name, durable, autodelete, exclusive, noWait, nil)

	if err != nil {
		return nil, err
	}

	return &Queue{&queue}, nil
}

// Publish publishes the given data under the given content type
func (c *Channel) Publish(queue *Queue, data string, contentType string) error {
	return c.Instance.Publish(
		"",
		queue.Instance.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(data),
		},
	)
}

// PublishAsJSON publishes the given data structure in json format
func (c *Channel) PublishAsJSON(queue *Queue, data interface{}) error {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	return c.Publish(queue, string(jsonData), "application/json")
}

// Close closes the opened channel
func (c *Channel) Close() {
	c.Instance.Close()
}
