package stm

import "bytes"

type BufferAdapter struct {
	buffers []*bytes.Buffer
}

func NewBufferAdapter() *BufferAdapter {
	return &BufferAdapter{}
}

func (a *BufferAdapter) Bytes() [][]byte {
	outBuffer := make([][]byte, len(a.buffers))

	for i, buf := range a.buffers {
		outBuffer[i] = buf.Bytes()
	}

	return outBuffer
}

func (a *BufferAdapter) Write(_ *Location, data []byte) {
	a.buffers = append(a.buffers, bytes.NewBuffer(data))
}
