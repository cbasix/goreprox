package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPacked(t *testing.T) {
	var buf bytes.Buffer
	frames := make(chan Frame, 1)

	frame := Frame{
		ConnectionId:   0,
		DropConnection: true,
		Data:           []byte{'i', 'o'},
	}
	frames <- frame

	go packedWriter(&buf, frames)
	go packedReader(&buf, frames)

	retrieved := <-frames

	if fmt.Sprintf("%+v", retrieved) != fmt.Sprintf("%+v", frame) {
		t.Errorf("Rounttrip changed something. Expected %+v got %+v", frame, retrieved)
	}
}
