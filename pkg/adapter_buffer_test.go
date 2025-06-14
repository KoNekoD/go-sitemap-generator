package stm

import "testing"

func TestBufferAdapter(t *testing.T) {
	a := NewBufferAdapter()

	a.Write(nil, []byte("Hello world"))

	a.Bytes()
}
