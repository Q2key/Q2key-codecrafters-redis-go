package core

import (
	"errors"
	"fmt"
	"os"
)

type RedisDB struct {
	Path    string
	Meta    map[string]string
	Aux     map[string]string
	Data    map[string]string
	Expires map[string]uint64
}

func NewRedisDB(path string) *RedisDB {
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

	r.Data = make(map[string]string)
	r.Expires = make(map[string]uint64)

	j := 0
	for i, b := range buff {
		if b == 0xFC {
			x := i + 1
			y := x + 8

			b := buff[x:y]
			exp := ParseMSecDateTimeStamp(&b)
			ok, key, _ := ParseValuePair(y+1, &buff)
			if ok {
				r.Expires[*key] = exp
			}

			j = y
		}

		if i < j {
			continue
		}

		if b == 0x00 {
			ok, key, val := ParseValuePair(i+1, &buff)
			if ok {
				r.Data[*key] = *val
			}
		}
	}

	fmt.Printf("\r\n%v", r.Expires)
	return nil
}

func (r *RedisDB) ReadFrom() (error, *Instance) {
	return nil, nil
}

func (r *RedisDB) Save(store *Instance) (error, *Instance) {
	return nil, nil
}

func (r *RedisDB) Flush(buff []byte, file os.File) error {
	return nil
}
