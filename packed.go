package main

import (
	"encoding/gob"
	"io"
)

/*
Reads frames from the given reader and puts them into the given channel.
Used for getting the frames out of the share connection and send them to further processing
*/
func packedReader(in io.Reader, frames chan<- Frame) {
	for {
		var frm Frame

		dec := gob.NewDecoder(in)
		err := dec.Decode(&frm)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		frames <- frm
	}
}

/*
Write frames from the chanel onto the output.
Used for coping the frames onto the network connections
*/
func packedWriter(out io.Writer, frames <-chan Frame) {
	for frm := range frames {
		enc := gob.NewEncoder(out)
		err := enc.Encode(frm)
		if err != nil {
			panic(err)
		}
	}
}
