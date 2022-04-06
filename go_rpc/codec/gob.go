package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	coon io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder
	enc  *gob.Encoder
}

var _gobcodec = (*JsonCodec)(nil)

func NewGobCodec(coon io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(coon)
	return &GobCodec{
		coon: coon,
		buf:  buf,
		dec:  gob.NewDecoder(coon),
		enc:  gob.NewEncoder(buf),
	}
}

func (j *GobCodec) ReadHeader(header *Header) error {
	return j.dec.Decode(header)
}

func (j *GobCodec) ReadBody(x interface{}) error {
	return j.dec.Decode(x)
}

func (j *GobCodec) Write(header *Header, x interface{}) (err error) {
	defer func() {
		_ = j.buf.Flush()
		if err != nil {
			_ = j.Close()
		}
	}()

	if err := j.enc.Encode(header); err != nil {
		log.Println("rpc codec: gob error encoding header:", err)
		return err
	}
	if err := j.enc.Encode(x); err != nil {
		log.Println("rpc codec: gob error encoding body:", err)
		return err
	}

	return nil
}

func (j *GobCodec) Close() error {
	return j.coon.Close()
}
