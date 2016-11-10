package utils

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestNewLimitlessReader(t *testing.T) {
	testcases := []struct {
		base     string
		expected *LimitlessReader
	}{
		{
			base:     "123",
			expected: &LimitlessReader{base: strings.NewReader("123")},
		},
		{
			base:     "",
			expected: &LimitlessReader{base: strings.NewReader(default_base_string)},
		},
	}

	for index, testcase := range testcases {
		actual := NewLimitlessReader(testcase.base)
		if reflect.DeepEqual(actual, testcase.expected) == false {
			t.Errorf("Testcase.%d: expected %v but %v\n",
				index,
				testcase.expected,
				actual,
			)
		}
	}
}

func TestLimitlessReaderRead1(t *testing.T) {
	reader := NewLimitlessReader("")

	buff := make([]byte, 128)
	for i := 0; i < 100; i++ {
		_, err := reader.Read(buff)
		if err != nil {
			t.Errorf("Fail to Read Data(%s)\n", err.Error())
			return
		}

	}
}

func TestLimitlessReaderRead2(t *testing.T) {
	reader := NewLimitlessReader("abc")

	buff := make([]byte, 128)
	for i := 0; i < 100; i++ {
		_, err := reader.Read(buff)
		if err != nil {
			t.Errorf("Fail to Read Data(%s)\n", err.Error())
			return
		}

		if bytes.HasPrefix(buff, []byte("abc")) == false {
			t.Errorf("Invalid data readed, buff should start with 'abc'")
		}
	}
}

func TestLimitlessReaderRead3(t *testing.T) {
	reader := NewLimitlessReader("")

	buff := make([]byte, 10240)
	for i := 0; i < 100; i++ {
		_, err := reader.Read(buff)
		if err != nil {
			t.Errorf("Fail to Read Data(%s)\n", err.Error())
			return
		}

	}
}
