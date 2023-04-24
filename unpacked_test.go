package main

import (
	"bytes"
	"testing"
)

func TestUnpacked(t *testing.T) {
	var buf bytes.Buffer
	frames := make(chan Frame, 1)

	frame := Frame{
		ConnectionId:   0,
		DropConnection: false,
		Data:           []byte{'i', 'o'},
	}
	frames <- frame

	go unpackedWriter(&buf, frames)
	go unpackedReader(&buf, frames)

	retrieved := <-frames

	if string(retrieved.Data) != string(frame.Data) {
		t.Errorf("Rounttrip changed something. Expected %+v got %+v", frame.Data, retrieved.Data)
	}
}
