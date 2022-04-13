package codec

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

type GobCodec struct {
	conn io.ReadWriteCloser //
	buf  *bufio.Writer
	dec  *gob.Decoder // 可以编码
	enc  *gob.Encoder // 可以解码
}

var _ Codec = (*GobCodec)(nil) // gob转换成通用的coder 这句是为了编译吧
// codec 是编解码器的意思. compress decompress

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn) // 用 bufio来初始化 4096个字节的缓冲区.
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h) // 这个是decode到传入的变量中.
}

func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil { // 这个是直接进行转换的吧.
		log.Println("rpc: gob error encoding header:", err)
		return
	}
	if err = c.enc.Encode(body); err != nil { // 把 body 也写入进去.
		log.Println("rpc: gob error encoding body:", err)
		return
	}
	return
}

func (c *GobCodec) Close() error {
	return c.conn.Close()
}
