package utils

type MessageQueue struct {
	queue chan string
}

type DLQueue struct {
	queue chan string
}

func (q *MessageQueue) Enqueue(msg string) {
	q.queue <- msg
}

func (q *MessageQueue) Dequeue() string {
	return <-q.queue
}

func (q *DLQueue) Enqueue(msg string) {
	q.queue <- msg
}

func (q *DLQueue) Dequeue() string {
	return <-q.queue
}
