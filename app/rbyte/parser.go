package rbyte

import (
	"encoding/binary"
	"time"
)

/*
Byte 	Name 			Description
0xFF 	EOF 			End of the RDB file
0xFE 	SELECTDB 		Database Selector
0xFD 	EXPIRETIME 		Expire time in seconds, see Key Expiry Timestamp
0xFC 	EXPIRETIMEMS 	Expire time in milliseconds, see Key Expiry Timestamp
0xFB 	RESIZEDB 		Hash table sizes for the main keyspace and expires, see Resizedb information
0xFA 	AUX 			Auxiliary fields. Arbitrary key-value settings, see Auxiliary fields
*/
const (
	AUX          = 0xFA
	RESIZEDB     = 0xFB
	EXPIRETIMEMS = 0xFC
	EXPIRETIME   = 0xFD
	SELECTDB     = 0xFE
	EOF          = 0xFF
	NULLISH_BYTE = 0b00
)

func ParseValuePair(i int, buff *[]byte) (bool, *string, *string) {
	sb := (*buff)[i:]
	first := sb[0]

	if !CheckByte(first) {
		return false, nil, nil
	}

	kl := int(first)

	vbf := sb[1+kl:]
	vb := vbf[0]
	if !CheckByte(vb) {
		return false, nil, nil
	}

	vl := int(vb)

	key := string(sb[1 : kl+1])
	val := string(vbf[1 : vl+1])

	return true, &key, &val
}

func GetDateFromTimeStamp(ts uint64) time.Time {
	return time.UnixMilli(int64(ts)).UTC()
}

func ParseMSecDateTimeStamp(buff *[]byte) uint64 {
	return binary.LittleEndian.Uint64(*buff)
}

func ParseSecDateTimeStamp(buff *[]byte) uint64 {
	return binary.LittleEndian.Uint64(*buff)
}

func CheckByte(b byte) bool {
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
	case NULLISH_BYTE:
		return false
	case EOF:
		return false
	}
	return true
}
