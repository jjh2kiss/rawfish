package math

import "testing"

func TestIntMin(t *testing.T) {
	testcases := []struct {
		a        int
		b        int
		expected int
	}{
		{1, 2, 1},
		{1, 1, 1},
		{2, 1, 1},
	}

	for index, testcase := range testcases {
		actual := IntMin(testcase.a, testcase.b)
		if actual != testcase.expected {
			t.Errorf("Testcase.%d: expected %d but %d\n",
				index,
				testcase.expected,
				actual,
			)
		}
	}
}

func TestIntMax(t *testing.T) {
	testcases := []struct {
		a        int
		b        int
		expected int
	}{
		{1, 2, 2},
		{1, 1, 1},
		{2, 1, 2},
	}

	for index, testcase := range testcases {
		actual := IntMax(testcase.a, testcase.b)
		if actual != testcase.expected {
			t.Errorf("Testcase.%d: expected %d but %d\n",
				index,
				testcase.expected,
				actual,
			)
		}
	}
}
