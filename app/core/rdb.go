package core

import (
	"errors"
	"fmt"
	"os"
)

type Rdb struct {
	Path    string
	Meta    map[string]string
	Aux     map[string]string
	data    map[string]string
	expires map[string]uint64
}

func NewRedisDB(path string) *Rdb {
	return &Rdb{
		Path: path,
	}
}

func (r *Rdb) Create() error {
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

func (r *Rdb) IsFileExists(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// Connect https://app.codecrafters.io/courses/redis/stages/jz6
// Connect https://rdb.fnordig.de/file_format.html
func (r *Rdb) Connect() error {
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
		if b == EOF {
			break
		}

		if b == RESIZEDB {
			x := i + 1

			bs := (buff[x] >> 6) & 0b00000011
			tb := fmt.Sprintf("%02b", bs)
			// The next 6 bits represent the length
			if tb == "00" {
				j = x + 1
			}

			// Read one additional byte. The combined 14 bits represent the length
			if tb == "01" {
				j = x + 2
			}

			// Discard the remaining 6 bits. The next 4 bytes from the stream represent the length
			if tb == "10" {
				j = x + 4
			}

			// The next object is encoded in a special format. The remaining 6 bits indicate the format.
			// May be used to store numbers or Strings, see String Encoding
			if tb == "11" {
				j = x + 2
			}
		}

		if b == EXPIRETIMEMS {
			x := i + 1
			y := x + 8

			sb := buff[x:y]
			exp := ParseMSecDateTimeStamp(&sb)
			ok, key, _ := ParseValuePair(y+1, &buff)
			if ok {
				r.expires[*key] = exp
			}

			j = y
		}

		if b == EXPIRETIME {
			x := i + 1
			y := x + 4
			sb := buff[x:y]
			exp := ParseSecDateTimeStamp(&sb)
			ok, key, _ := ParseValuePair(y+1, &buff)
			if ok {
				r.expires[*key] = exp
			}

			j = y
		}

		if i < j {
			continue
		}

		if b == 0x00 {
			ok, key, val := ParseValuePair(i+1, &buff)
			if ok {
				r.data[*key] = *val
			}
		}
	}

	return nil
}

func (r *Rdb) GetData() map[string]string {
	return r.data
}

func (r *Rdb) GetExpires() map[string]uint64 {
	return r.expires
}
