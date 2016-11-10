package service

import "testing"

func TestTypeIsNormalType(t *testing.T) {
	testcases := []struct {
		in       int
		expected bool
	}{
		{in: SERVICETYPE_NORMAL, expected: true},
		{in: SERVICETYPE_RAW, expected: false},
	}

	for index, testcase := range testcases {
		actual := Type(testcase.in).IsNormalType()
		if actual != testcase.expected {
			t.Errorf("Testcase.%d: expected %v but %v\n",
				index,
				testcase.expected,
				actual)
		}
	}
}

func TestTypeIsRawType(t *testing.T) {
	testcases := []struct {
		in       int
		expected bool
	}{
		{in: SERVICETYPE_NORMAL, expected: false},
		{in: SERVICETYPE_RAW, expected: true},
	}

	for index, testcase := range testcases {
		actual := Type(testcase.in).IsRawType()
		if actual != testcase.expected {
			t.Errorf("Testcase.%d: expected %v but %v\n",
				index,
				testcase.expected,
				actual)
		}
	}
}
