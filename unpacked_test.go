package main

import (
	"bytes"
	"fmt"
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

	fmt.Print(string(frame.Data))
	if string(retrieved.Data) != string(frame.Data) {
		t.Errorf("Rounttrip changed something. Expected %+v got %+v", frame.Data, retrieved.Data)
	}
}

func TestUnpackedWriterRaw(t *testing.T) {
	var buf bytes.Buffer
	out := make(chan Frame, 1)
	out <- Frame{Data: []byte{1, 2, 3}}
	close(out)

	unpackedWriter(&buf, out)

	result := make([]byte, 3)
	buf.Read(result)

	if result[0] != 1 || result[2] != 3 {
		t.Errorf("Raw write produced wrote invalid bytes expected: 1, 2, 3 got: %+v", result)
	}
}

func TestUnpackedReaderRaw(t *testing.T) {
	var buf bytes.Buffer
	in := make(chan Frame, 1)
	buf.Write([]byte{1, 2, 3})

	unpackedReader(&buf, in)

	frm := <-in

	if frm.Data[0] != 1 || frm.Data[2] != 3 {
		t.Errorf("Raw read produced frame with invalid data expected data: 1, 2, 3 got: %+v", frm)
	}
}
