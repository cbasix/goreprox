package main

import (
	"encoding/binary"
	"encoding/gob"
	"io"
)

/*
Reads frames from the given reader and puts them into the given channel.
Used for getting the frames out of the share connection and send them to further processing
*/
func packedGobReader(in io.Reader, frames chan<- Frame) {
	for {
		var frm Frame

		dec := gob.NewDecoder(in)
		err := dec.Decode(&frm)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		logD("packed send: %+v", frm)
		frames <- frm
	}
}

/*
Write frames from the chanel onto the output.
Used for coping the frames onto the network connections
*/
func packedGobWriter(out io.Writer, frames <-chan Frame) {
	for frm := range frames {
		logD("packed send: %+v", frm)
		enc := gob.NewEncoder(out)
		err := enc.Encode(frm)
		if err != nil {
			panic(err)
		}
	}
}

/*
Reads frames from the given reader and puts them into the given channel.
Used for getting the frames out of the share connection and send them to further processing
*/
func packedReader(in io.Reader, frames chan<- Frame) {
	for {
		var frm Frame
		header := make([]byte, 4)
		read, err := in.Read(header)
		logD("Read header %v", header)
		if err == io.EOF {
			return
		} else if read < 4 {
			panic("Short read header")
		} else if err != nil {
			panic(err)
		}

		packageLen := uint(binary.BigEndian.Uint16(header[2:]))
		if packageLen > 4086 {
			panic("package length to big")
		}
		buf := make([]byte, packageLen)
		read, err = in.Read(buf)
		if err == io.EOF {
			close(frames)
			return
		} else if read < int(packageLen) {
			panic("Short read data")
		} else if err != nil {
			panic(err)
		}

		frm = Frame{
			ConnectionId:   uint(header[0]),
			DropConnection: header[1] != 0,
			Data:           buf,
		}

		logD("packed read: %+v", frm)
		frames <- frm
	}
}

/*
Write frames from the chanel onto the output.
Used for coping the frames onto the network connections
*/
func packedWriter(out io.Writer, frames <-chan Frame) {
	for frm := range frames {
		var header uint32 = 0
		header |= uint32(frm.ConnectionId) << 24
		if frm.DropConnection {
			header |= 1 << 16
		}
		header |= uint32(len(frm.Data))

		headBuff := make([]byte, 4)
		binary.BigEndian.PutUint32(headBuff, header)

		logD("Sending---\nHeader: %v", headBuff)
		written, err := out.Write(headBuff)
		if written < 4 {
			panic("Short write header")
		} else if err != nil {
			panic(err)
		}
		logD("Data %v", frm.Data)
		written, err = out.Write(frm.Data)
		if written < int(len(frm.Data)) {
			panic("Short write data")
		} else if err != nil {
			panic(err)
		}

		logD("packed send: %+v", frm)

	}
}
