package codec

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

type JsonCodec struct {
	coon io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder
	enc  *json.Encoder
}

var _Codec = (*JsonCodec)(nil)

func NewJsonCodec(coon io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(coon)
	return &JsonCodec{
		coon: coon,
		buf:  buf,
		dec:  json.NewDecoder(coon),
		enc:  json.NewEncoder(buf),
	}
}

func (j *JsonCodec) ReadHeader(header *Header) error {
	return j.dec.Decode(header)
}

func (j *JsonCodec) ReadBody(x interface{}) error {
	return j.dec.Decode(x)
}

func (j *JsonCodec) Write(header *Header, x interface{}) (err error) {
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

func (j *JsonCodec) Close() error {
	return j.coon.Close()
}
