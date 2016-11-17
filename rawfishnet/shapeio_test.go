package rawfishnet

import (
	"bufio"
	"bytes"
	"math"
	"testing"
	"time"
)

type TmpWriter struct {
	*bufio.Writer
}

func (self TmpWriter) Flush() {
	_ = self.Writer.Flush()
}

func TestWriteFlushWrite(t *testing.T) {
	b := bytes.NewBuffer(nil)

	writer := TmpWriter{bufio.NewWriter(b)}

	wf := WriteFlush{writer, writer}

	wf.Write([]byte("hello world"))

	if b.String() != "hello world" {
		t.Errorf("Fail to write and flush")
	}
}

func TestWriteFlushWithoutFlush(t *testing.T) {
	b := bytes.NewBuffer(nil)

	writer := TmpWriter{bufio.NewWriter(b)}

	writer.Write([]byte("hello world"))

	if b.String() == "hello world" {
		t.Errorf("Fail to write and flush")
	}
}

func TestCopyWithShapeIO(t *testing.T) {
	testcases := []struct {
		input    []byte
		rate     int
		expected float64
	}{
		{
			input:    []byte("123"),
			rate:     1,   //1byte
			expected: 3.0, //3seconds
		},
		{
			input:    []byte("1234567890"),
			rate:     2,   //1byte
			expected: 5.0, //5seconds
		},
		{
			input:    bytes.Repeat([]byte("1234567890"), 5),
			rate:     10,  //10byte/second
			expected: 5.0, //5seconds
		},
	}

	for index, testcase := range testcases {
		src := bytes.NewReader(testcase.input)
		dst := TmpWriter{bufio.NewWriter(nil)}

		begin := time.Now()
		written, err := CopyWithShapeIO(dst, src, len(testcase.input), testcase.rate)
		duration := time.Since(begin).Seconds()

		if err != nil {
			t.Errorf("Testcase.%d : %s", index, err.Error())
			return
		}
		if written != len(testcase.input) {
			t.Errorf("Testcase.%d : expected written lenght is %d, but %d", index, testcase.expected, written)
			return
		}

		difference := math.Abs(duration - testcase.expected)
		//오차가 1보다 크다면
		if difference >= 1 {
			t.Errorf("Testcase.%d : expected duration time is %f, but %f", index, testcase.expected, duration)
		}

	}
}
