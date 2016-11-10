package service

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/golang/groupcache/lru"
)

func TestNewServiceType(t *testing.T) {
	curr, err := os.Getwd()
	if err != nil {
		t.Errorf("Fail to get Current Working Directory")
		return
	}

	testcases := []struct {
		root         string
		default_type Type
		expected     *ServiceType
	}{
		{
			root:         "",
			default_type: Type(SERVICETYPE_NORMAL),
			expected:     nil,
		},
		{
			root:         "/",
			default_type: Type(SERVICETYPE_NORMAL),
			expected: &ServiceType{
				root:         "/",
				default_type: Type(SERVICETYPE_NORMAL),
				cache:        lru.New(LRU_CACHE_SIZE),
			},
		},
		{
			root:         "./",
			default_type: Type(SERVICETYPE_NORMAL),
			expected: &ServiceType{
				root:         curr,
				default_type: Type(SERVICETYPE_NORMAL),
				cache:        lru.New(LRU_CACHE_SIZE),
			},
		},
	}

	for index, testcase := range testcases {
		actual := NewServiceType(testcase.root, testcase.default_type)
		if reflect.DeepEqual(actual, testcase.expected) == false {
			t.Errorf("Testcase.%d:\nexpected\n%v\nbut\n%v\n",
				index,
				testcase.expected,
				actual)
		}
	}
}

func TestInternalGet(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Fail to get Current Working Directory")
		return
	}

	testcases := []struct {
		root         string
		default_type Type
		path         string
		expected     Type
	}{
		{
			root:         filepath.Join(cwd, "test"),
			default_type: Type(SERVICETYPE_NORMAL),
			path:         "normal",
			expected:     Type(SERVICETYPE_NORMAL),
		},
		{
			root:         filepath.Join(cwd, "test"),
			default_type: Type(SERVICETYPE_NORMAL),
			path:         "raw",
			expected:     Type(SERVICETYPE_RAW),
		},
		{
			root:         filepath.Join(cwd, "test"),
			default_type: Type(SERVICETYPE_RAW),
			path:         "normal",
			expected:     Type(SERVICETYPE_RAW),
		},
		//test/raw/.raw의 영향으로 raw 타입
		{
			root:         filepath.Join(cwd, "test"),
			default_type: Type(SERVICETYPE_NORMAL),
			path:         "raw/sub1/sub2",
			expected:     Type(SERVICETYPE_RAW),
		},
		//test/raw/.raw의 영향으로 raw 타입
		{
			root:         filepath.Join(cwd, "test"),
			default_type: Type(SERVICETYPE_RAW),
			path:         "raw/sub1/sub2",
			expected:     Type(SERVICETYPE_RAW),
		},
		{
			root:         filepath.Join(cwd, "test"),
			default_type: Type(SERVICETYPE_NORMAL),
			path:         "raw/sub1/sub2/a.txt",
			expected:     Type(SERVICETYPE_RAW),
		},
	}

	for index, testcase := range testcases {
		servicetype := NewServiceType(testcase.root, testcase.default_type)
		actual := servicetype.get(testcase.path)
		if actual != testcase.expected {
			t.Errorf("Testcase.%d:\nexpected\n%v\nbut\n%v\n",
				index,
				testcase.expected,
				actual)
		}
	}
}

func TestGet(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Fail to get Current Working Directory")
		return
	}

	testcases := []struct {
		root         string
		default_type Type
		path         string
		expected     Type
	}{
		{
			root:         filepath.Join(cwd, "test"),
			default_type: Type(SERVICETYPE_NORMAL),
			path:         "normal",
			expected:     Type(SERVICETYPE_NORMAL),
		},
	}
	for index, testcase := range testcases {
		servicetype := NewServiceType(testcase.root, testcase.default_type)
		actual := servicetype.Get(testcase.path)
		if actual != testcase.expected {
			t.Errorf("Testcase.%d:\nexpected\n%v\nbut\n%v\n",
				index,
				testcase.expected,
				actual)
		}
	}
}

func TestGetCache(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Fail to get Current Working Directory")
		return
	}

	servicetype := NewServiceType(
		filepath.Join(cwd, "test"),
		Type(SERVICETYPE_NORMAL))

	if servicetype.cache.Len() != 0 {
		t.Errorf("Wrong Cache Size, expected 0, but %d\n",
			servicetype.cache.Len(),
		)
		return
	}

	_ = servicetype.Get("normal")

	if servicetype.cache.Len() != 1 {
		t.Errorf("Wrong Cache Size, expected 1, but %d\n",
			servicetype.cache.Len(),
		)
	}
}
