package utils

import (
	"testing"
	"time"
)

func TestMessageQueueInit(t *testing.T) {
	q := &MessageQueue{}
	q.Init(4)

	if cap(q.queue) != 4 {
		t.Errorf("expected cap 4, got %d", cap(q.queue))
	}
	if len(q.queue) != 0 {
		t.Errorf("expected len 0, got %d", len(q.queue))
	}
}

func TestMessageQueueEnqueueDequeue(t *testing.T) {
	q := &MessageQueue{}
	q.Init(8)

	q.Enqueue("alpha")
	q.Enqueue("beta")
	q.Enqueue("gamma")

	if v := q.Dequeue(); v != "alpha" {
		t.Errorf("expected alpha, got %s", v)
	}
	if v := q.Dequeue(); v != "beta" {
		t.Errorf("expected beta, got %s", v)
	}
	if v := q.Dequeue(); v != "gamma" {
		t.Errorf("expected gamma, got %s", v)
	}
}

func TestMessageQueueGetChannel(t *testing.T) {
	q := &MessageQueue{}
	q.Init(4)

	q.Enqueue("a")
	q.Enqueue("b")

	ch := q.GetChannel()

	select {
	case v := <-ch:
		if v != "a" {
			t.Errorf("expected a, got %s", v)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for first message")
	}

	select {
	case v := <-ch:
		if v != "b" {
			t.Errorf("expected b, got %s", v)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for second message")
	}
}

func TestMessageQueueBufferedEnqueueNonBlocking(t *testing.T) {
	q := &MessageQueue{}
	q.Init(3)

	q.Enqueue("x")
	q.Enqueue("y")
	q.Enqueue("z")

	if len(q.queue) != 3 {
		t.Errorf("expected len 3, got %d", len(q.queue))
	}
}

func TestDLQueueInit(t *testing.T) {
	q := &DLQueue{}
	q.Init(8)

	if cap(q.queue) != 8 {
		t.Errorf("expected cap 8, got %d", cap(q.queue))
	}
}

func TestDLQueueEnqueueDequeue(t *testing.T) {
	q := &DLQueue{}
	q.Init(4)

	q.Enqueue("dead1")
	q.Enqueue("dead2")

	if v := q.Dequeue(); v != "dead1" {
		t.Errorf("expected dead1, got %s", v)
	}
	if v := q.Dequeue(); v != "dead2" {
		t.Errorf("expected dead2, got %s", v)
	}
}

func TestDLQueueGetChannel(t *testing.T) {
	q := &DLQueue{}
	q.Init(4)

	q.Enqueue("dlq-a")
	ch := q.GetChannel()

	select {
	case v := <-ch:
		if v != "dlq-a" {
			t.Errorf("expected dlq-a, got %s", v)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for dlq message")
	}
}
