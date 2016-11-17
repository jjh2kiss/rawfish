package rawfishnet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/fujiwara/shapeio"
	"github.com/jjh2kiss/rawfish/math"
)

//ShapeIOResponseWriter implements http.ResponseWriter
type ShapeIOResponseWriter struct {
	base    http.ResponseWriter
	writer  *shapeio.Writer
	flusher http.Flusher
	rate    int
}

func NewShapeIOResponseWriter(base http.ResponseWriter, rate int) *ShapeIOResponseWriter {
	writer := shapeio.NewWriter(base)
	writer.SetRateLimit(float64(rate))
	flusher := base.(http.Flusher)

	return &ShapeIOResponseWriter{
		base:    base,
		writer:  writer,
		flusher: flusher,
		rate:    rate,
	}
}

func (self *ShapeIOResponseWriter) Write(b []byte) (written int, err error) {
	src := bytes.NewReader(b)
	dst := WriteFlush{self.writer, self.flusher}

	return CopyWithShapeIO(dst, src, len(b), self.rate)
}

func (self *ShapeIOResponseWriter) Header() http.Header {
	return self.base.Header()
}

func (self *ShapeIOResponseWriter) WriteHeader(s int) {
	self.base.WriteHeader(s)
}

func (self *ShapeIOResponseWriter) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	hj, ok := self.base.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("Fail to assert Hijacker")
	}

	return hj.Hijack()
}

type WriteFlush struct {
	writer  io.Writer
	flusher http.Flusher
}

func (self WriteFlush) Write(b []byte) (written int, err error) {
	if self.writer != nil {
		written, err = self.writer.Write(b)
		if self.flusher != nil {
			self.flusher.Flush()
		}
		return
	}
	return 0, fmt.Errorf("Empty writer")
}

func CopyWithShapeIO(dst io.Writer, src io.Reader, size int, rate int) (written int, err error) {
	buf_size := math.IntMin(size, rate)
	buf := make([]byte, buf_size)

	writer := shapeio.NewWriter(dst)
	writer.SetRateLimit(float64(rate))

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := writer.Write(buf[0:nr])

			if nw > 0 {
				written += nw
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err

}
