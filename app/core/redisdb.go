package core

import (
	"errors"
	"os"
)

type RedisDB struct {
	Path string
	Meta map[string]string
	Aux  map[string]string
	Data map[string]string
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
		} else if b == EXPIRETIMEMS {
			//fmt.Printf("Expiretime ms")
		} else if b == EXPIRETIME {
			//fmt.Printf("Expiretime")
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

func (r *RedisDB) Flush(buff []byte, file os.File) error {
	return nil
}

const (
	AUX          = 0xFA
	RESIZEDB     = 0xFB
	EXPIRETIMEMS = 0xFC
	EXPIRETIME   = 0xFD
	SELECTDB     = 0xFE
	EOF          = 0xFF
)

const NULL = 0b00

/*
Byte 	Name 	Description
0xFF 	EOF 	End of the RDB file
0xFE 	SELECTDB 	Database Selector
0xFD 	EXPIRETIME 	Expire time in seconds, see Key Expiry Timestamp
0xFC 	EXPIRETIMEMS 	Expire time in milliseconds, see Key Expiry Timestamp
0xFB 	RESIZEDB 	Hash table sizes for the main keyspace and expires, see Resizedb information
0xFA 	AUX 	Auxiliary fields. Arbitrary key-value settings, see Auxiliary fields
*/

func (r *RedisDB) checkByte(b byte) bool {
	switch b {
	case AUX:
		return false
	case SELECTDB:
		return false
	case RESIZEDB:
		return false
	case EXPIRETIME:
		return false
	case EXPIRETIMEMS:
		return false
	case NULL:
		return false
	case EOF:
		return false
	}
	return true
}
