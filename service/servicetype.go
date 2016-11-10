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

func (self *ServiceType) get(path string) Type {
	path = filepath.Clean(path)

	filename := filepath.Join(self.root, path, ".raw")
	filename = filepath.Clean(filename)

	_, err := os.Stat(filename)

	if err != nil {
		//Root 디렉토리 검사 후 파일이 존재하지 않는다면
		//Default type을 리턴한다.
		if path == "." {
			return self.default_type
		}

		path = filepath.Dir(path)
		return self.get(path)
	}

	//stat이 있다면, .raw 라는 이름을 가진
	//dir, file, link 가 있다는 것이다.
	//Raw를 리턴하자.
	return Type(SERVICETYPE_RAW)
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
