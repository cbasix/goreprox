package main

import (
	"io"
)

/*
Reads frames from the given reader and puts them into the given channel.
Used for getting the frames out of the share connection and send them to further processing
*/
func unpackedReader(in io.Reader, frames chan<- Frame) {
	buf := make([]byte, 4096)
	for {
		read, err := in.Read(buf)

		frames <- Frame{
			ConnectionId:   0, // will be set later on
			DropConnection: err == io.EOF,
			Data:           buf[:read],
		}
	}
}

/*
Write frames from the chanel onto the output.
Used for coping the frames onto the network connections
*/
func unpackedWriter(out io.Writer, frames <-chan Frame) {
	for frm := range frames {
		_, err := out.Write(frm.Data)
		if err != nil {
			panic(err)
		}
	}
}
