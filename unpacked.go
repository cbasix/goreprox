package main

import (
	"errors"
	"io"
	"net"
)

/*
Reads frames from the given reader and puts them into the given channel.
Used for getting the frames out of the share connection and send them to further processing
*/
func unpackedReader(in io.Reader, frames chan<- Frame) {
	buf := make([]byte, 4096)
	for {
		read, err := in.Read(buf)

		if err != nil && err != io.EOF {
			if errors.Is(err, net.ErrClosed) {
				close(frames)
				return
			}
			panic(err)
		}
		if read > 0 {
			frm := Frame{
				ConnectionId:   0, // will be set later on
				DropConnection: err == io.EOF,
				Data:           buf[:read],
			}
			frames <- frm
			//logD("unpacked receive: %+v", frm)
		}

		if err == io.EOF {
			close(frames)
			break
		}
	}
}

/*
Write frames from the chanel onto the output.
Used for coping the frames onto the network connections
*/
func unpackedWriter(out io.Writer, frames <-chan Frame) {
	for frm := range frames {
		logD("unpacked send: %+v", frm)
		_, err := out.Write(frm.Data)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			panic(err)
		}
	}
}
