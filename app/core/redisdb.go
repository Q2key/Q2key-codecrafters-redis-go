package core

import (
	"errors"
	"github.com/codecrafters-io/redis-starter-go/app/contracts"
	"github.com/codecrafters-io/redis-starter-go/app/rbyte"
	"os"
)

type RedisDB struct {
	Path    string
	Meta    map[string]string
	Aux     map[string]string
	data    map[string]string
	expires map[string]uint64
}

func NewRedisDB(path string) contracts.Database {
	return &RedisDB{
		Path: path,
	}
}

func (r *RedisDB) Create() error {
	if r.IsFileExists(r.Path) {
		return errors.New("file exists")
	}

	f, err := os.Create(r.Path)
	defer f.Close()

	if err != nil {
		return err
	}

	return nil
}

func (r *RedisDB) IsFileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func (r *RedisDB) Connect() error {
	if !r.IsFileExists(r.Path) {
		return errors.New("file not exists")
	}

	f, err := os.Open(r.Path)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err != nil {
		return err
	}

	s, _ := f.Stat()

	buff := make([]byte, s.Size())

	_, err = f.Read(buff)
	if err != nil {
		return err
	}

	r.data = make(map[string]string)
	r.expires = make(map[string]uint64)

	j := 0
	for i, b := range buff {
		if b == rbyte.EOF {
			break
		}

		if b == rbyte.EXPIRETIMEMS {
			x := i + 1
			y := x + 8

			sb := buff[x:y]
			exp := rbyte.ParseMSecDateTimeStamp(&sb)
			ok, key, _ := rbyte.ParseValuePair(y+1, &buff)
			if ok {
				r.expires[*key] = exp
			}

			j = y
		}

		if b == rbyte.EXPIRETIME {
			x := i + 1
			y := x + 4
			sb := buff[x:y]
			exp := rbyte.ParseSecDateTimeStamp(&sb)
			ok, key, _ := rbyte.ParseValuePair(y+1, &buff)
			if ok {
				r.expires[*key] = exp
			}

			j = y
		}

		if i < j {
			continue
		}

		if b == 0x00 {
			ok, key, val := rbyte.ParseValuePair(i+1, &buff)
			if ok {
				r.data[*key] = *val
			}
		}
	}

	return nil
}

func (r *RedisDB) Data() map[string]string {
	return r.data
}
func (r *RedisDB) Expires() map[string]uint64 {
	return r.expires
}
