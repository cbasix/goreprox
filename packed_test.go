package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPacked(t *testing.T) {
	var buf bytes.Buffer
	out := make(chan Frame, 1)
	in := make(chan Frame, 1)

	frame := Frame{
		ConnectionId:   0,
		DropConnection: true,
		Data:           []byte{'i', 'o'},
	}
	out <- frame
	close(out)

	packedWriter(&buf, out)
	packedReader(&buf, in)

	retrieved := <-in

	if fmt.Sprintf("%+v", retrieved) != fmt.Sprintf("%+v", frame) {
		t.Errorf("Rounttrip changed something. Expected %+v got %+v", frame, retrieved)
	}
}
