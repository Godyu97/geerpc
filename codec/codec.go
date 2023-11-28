package codec

import (
	"io"
)

type Header struct {
	ServiceMethod string //format "Service.Method"
	Seq           uint64 //req id
	Error         string
}

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(any) error
	Write(*Header, any) error
}

type NewCodecFunc func(closer io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
	NewCodecFuncMap[JsonType] = NewJsonCodec
}

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber int  // MagicNumber marks this's a geerpc request
	CodecType   Type // client may choose different Codec to encode body
}
