package main

import (
	"bytes"
	"testing"
)

func TestUnpacked(t *testing.T) {
	var buf bytes.Buffer
	out := make(chan Frame, 1)
	in := make(chan Frame, 1)

	frame := Frame{
		ConnectionId:   0,
		DropConnection: false,
		Data:           []byte{'i', 'o'},
	}
	out <- frame
	close(out)

	unpackedWriter(&buf, out)
	unpackedReader(&buf, in)

	retrieved := <-in

	if string(retrieved.Data) != string(frame.Data) {
		t.Errorf("Rounttrip changed something. Expected %+v got %+v", frame.Data, retrieved.Data)
	}
}
