package frame

import (
	"encoding/binary"
	"io"
)

type FramePayload []byte

type StreamFrameCodec interface {
	Encode(writer io.Writer, payload FramePayload) error
	Decode(reader io.Reader) (FramePayload, error)
}

type myFrameCodec struct{}

func NewMyFrameCodec() StreamFrameCodec {
	return &myFrameCodec{}
}

func (p *myFrameCodec) Encode(w io.Writer, framePayload FramePayload) error {
	var f = framePayload
	var totalLen int32 = int32(len(framePayload)) + 4

	err := binary.Write(w, binary.BigEndian, &totalLen)
	if err != nil {
		return err
	}

	// make sure all data will be written to outbound stream
	for {
		n, err := w.Write(f)
		if err != nil {
			return err
		}
		if n >= len(f) {
			break
		}
		if n < len(f) {
			f = f[n:]
		}
	}
	return nil
}

func (p *myFrameCodec) Decode(r io.Reader) (FramePayload, error) {
	var totalLen int32
	err := binary.Read(r, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
