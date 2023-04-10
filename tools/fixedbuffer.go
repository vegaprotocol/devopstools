package tools

import (
	"sync"
)

type FixedBuffer struct {
	mutex sync.RWMutex
	data  []byte
	cap   int
}

func NewFixedBuffer(data []byte) FixedBuffer {
	return FixedBuffer{
		data: data,
		cap:  cap(data),
	}
}

func (b *FixedBuffer) Write(data []byte) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(data) > b.cap {
		b.data = data[len(data)-b.cap:]
		return
	}

	if b.cap < len(data)+len(b.data) {
		b.data = b.data[len(data):]
	}

	b.data = append(b.data, data...)
}

func (b *FixedBuffer) String() string {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return string(b.data)
}

func (b *FixedBuffer) Empty() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return len(b.data) == 0
}

func (b *FixedBuffer) Clear() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.data = []byte{}
}

func (b *FixedBuffer) Read() []byte {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	var data = make([]byte, 0, b.cap)
	copy(data, b.data)

	return data
}
