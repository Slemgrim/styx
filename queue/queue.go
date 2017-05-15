package queue

import (
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
func (q *Connection) Channel() {
	instance := q.Connection.Channel()

	return &Channel{instance}
}

// Close closes the opened channel
func (c *Channel) Close() {
	c.Instance.Close()
}
