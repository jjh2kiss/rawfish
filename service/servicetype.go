package service

import (
	"github.com/golang/groupcache/lru"

	"os"
	"path/filepath"
)

type ServiceType struct {
	root         string
	default_type Type
	cache        *lru.Cache
}

const LRU_CACHE_SIZE = 1024

func NewServiceType(root string, default_type Type) *ServiceType {
	if root == "" {
		return nil
	}

	root, err := filepath.Abs(root)
	if err != nil {
		return nil
	}

	return &ServiceType{
		root:         root,
		default_type: default_type,
		cache:        lru.New(LRU_CACHE_SIZE),
	}
}

func (self *ServiceType) exist(path string) bool {
	filename := filepath.Join(self.root, path)
	filename = filepath.Clean(filename)
	_, err := os.Stat(filename)

	if err != nil {
		return false
	}

	return true
}

func (self *ServiceType) get(path string) Type {
	path = filepath.Clean(path)

	if self.exist(filepath.Join(path, ".raw")) {
		return Type(SERVICETYPE_RAW)
	} else if self.exist(filepath.Join(path, ".normal")) {
		return Type(SERVICETYPE_NORMAL)
	}

	//Root 디렉토리 검사 후 파일이 존재하지 않는다면
	//Default type을 리턴한다.
	if path == "." || path == "/" {
		return self.default_type
	}

	path = filepath.Dir(path)
	return self.get(path)
}

//path is relpath
func (self *ServiceType) Get(path string) Type {
	var result Type

	value, ok := self.cache.Get(path)
	if ok == true {
		result = value.(Type)
	} else {
		result = self.get(path)
		self.cache.Add(path, result)
	}

	return result
}
