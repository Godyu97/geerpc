package codec

import (
	"bufio"
	"io"
	"github.com/Godyu97/geerpc/logger"
	jsoniter "github.com/json-iterator/go"
)

type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *jsoniter.Decoder
	enc  *jsoniter.Encoder
}

func (j *JsonCodec) Close() error {
	return j.conn.Close()
}

func (j *JsonCodec) ReadHeader(header *Header) error {
	return j.dec.Decode(header)
}

func (j *JsonCodec) ReadBody(body any) error {
	logger.Info("begin ReadBody json")
	return j.dec.Decode(body)
}

func (j *JsonCodec) Write(header *Header, body any) (err error) {
	defer func() {
		_ = j.buf.Flush()
		if err != nil {
			_ = j.Close()
		}
	}()
	if err = j.enc.Encode(header); err != nil {
		logger.Error("rpc codec: json error encoding header:", err)
		return err
	}
	if err = j.enc.Encode(body); err != nil {
		logger.Error("rpc codec: json error encoding body:", err)
		return err
	}
	logger.Info("write done:", *header, body)
	return nil
}

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  jsoniter.NewDecoder(conn),
		enc:  jsoniter.NewEncoder(buf),
	}
}
