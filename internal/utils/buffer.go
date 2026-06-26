package utils

import (
	"sync"

	"github.com/om-baji/write-go/internal"
)

type MemoryBuffer struct {
	mu   sync.Mutex
	data []byte
	size int64
}

func NewBuffer(capacity int) *MemoryBuffer {
	return &MemoryBuffer{
		data: make([]byte, 0, capacity),
		size: 0,
	}
}

func CommitWorker(message []byte, buffer *MemoryBuffer, segment internal.Segment, threshold int) (internal.Segment, error) {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()

	buffer.data = append(buffer.data, message...)
	buffer.data = append(buffer.data, '\n')
	buffer.size = int64(len(buffer.data))

	if buffer.size >= int64(threshold) {
		seg, err := AppendBuffer(segment, buffer.data)
		if err != nil {
			return segment, err
		}
		buffer.data = make([]byte, 0, threshold)
		buffer.size = 0
		return seg, nil
	}

	return segment, nil
}

func FlushBuffer(buffer *MemoryBuffer, segment internal.Segment) (internal.Segment, error) {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()

	if buffer.size == 0 {
		return segment, nil
	}

	seg, err := AppendBuffer(segment, buffer.data)
	if err != nil {
		return segment, err
	}

	buffer.data = make([]byte, 0, cap(buffer.data))
	buffer.size = 0
	return seg, nil
}
