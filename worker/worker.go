package worker

// QueueWorker defines the queue specific information
type QueueWorker struct {
	QueueName string
}

// NewQueueWorker creates a new queue worker instance
func NewQueueWorker(queueName string) *QueueWorker {
	return &QueueWorker{queueName}
}

// Start starts the worker execution
func (worker *QueueWorker) Start() {
	go func() {
		// tba
	}()
}
