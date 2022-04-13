package codec

import (
	"io"
)

type Header struct { // 用作消息的头.
	ServiceMethod string // format "Service.Method"
	Seq           uint64 // sequence number chosen by client
	Error         string
}

type Codec interface {
	io.Closer                         // 必须有 Close() 方法
	ReadHeader(*Header) error         // 读取到 header 中.
	ReadBody(interface{}) error       // 读取到 interface 中.
	Write(*Header, interface{}) error // 可以吧 header 和 body 写入进去.
}

type NewCodecFunc func(io.ReadWriteCloser) Codec // 可以讲一个转换过来吧.

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json" // not implemented
)

var NewCodecFuncMap map[Type]NewCodecFunc // 适用什么转换函数? 把 RWC 转换过来.

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
