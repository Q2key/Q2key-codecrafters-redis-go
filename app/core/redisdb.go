package core

import (
	"errors"
	"os"
)

type RedisDB struct {
	name string
	path string
	Meta map[string]string
	Aux  map[string]string
	Data map[string]string
}

func NewRedisDB(name string, path string) *RedisDB {
	return &RedisDB{
		name: name,
		path: path,
	}
}

func (r *RedisDB) Create(path string) error {
	if r.IsFileExists(path) {
		return errors.New("file exists")
	}

	f, err := os.Create(path)
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

func (r *RedisDB) Connect(path string) error {
	if !r.IsFileExists(path) {
		return errors.New("file not exists")
	}

	f, err := os.Open(path)
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

	r.Meta = make(map[string]string)
	r.Aux = make(map[string]string)
	r.Data = make(map[string]string)
	for i, b := range buff {
		if b == 0x00 {
			kbf := buff[i+1:]
			kb := kbf[0]
			if !r.checkByte(kb) {
				continue
			}

			kl := int(kb)

			vbf := kbf[1+kl:]
			vb := vbf[0]
			if !r.checkByte(vb) {
				continue
			}

			vl := int(vb)

			key := string(kbf[1 : kl+1])
			val := string(vbf[1 : vl+1])

			r.Data[key] = val
		}
	}

	return nil
}

func (r *RedisDB) ReadFrom() (error, *Instance) {
	return nil, nil
}

func (r *RedisDB) Save(store *Instance) (error, *Instance) {
	return nil, nil
}

func (r *RedisDB) Name() string {
	return r.name
}

func (r *RedisDB) Path() string {
	return r.path
}

func (r *RedisDB) Flush(buff []byte, file os.File) error {
	return nil
}

func (r *RedisDB) checkByte(b byte) bool {
	switch b {
	case 0xFA:
		return false
	case 0xFE:
		return false
	case 0xFB:
		return false
	case 0xFD:
		return false
	case 0xFC:
		return false
	case 0b00:
		return false
	default:
		return true
	}
}
