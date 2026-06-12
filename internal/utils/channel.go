package utils

type MessageQueue struct {
	queue chan string
}

type DLQueue struct {
	queue chan string
}

func (q *MessageQueue) Init(size int) {
	q.queue = make(chan string, size)
}

func (q *MessageQueue) Enqueue(msg string) {
	q.queue <- msg
}

func (q *MessageQueue) Dequeue() string {
	return <-q.queue
}

func (q *MessageQueue) GetChannel() <-chan string {
	return q.queue
}

func (q *DLQueue) Init(size int) {
	q.queue = make(chan string, size)
}

func (q *DLQueue) Enqueue(msg string) {
	q.queue <- msg
}

func (q *DLQueue) Dequeue() string {
	return <-q.queue
}

func (q *DLQueue) GetChannel() <-chan string {
	return q.queue
}
