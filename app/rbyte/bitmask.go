package rbyte

type bit byte

const (
	b0 bit = 0b0000_0001
	b1 bit = 0b0000_0010
	b2 bit = 0b0000_0100
	b3 bit = 0b0000_1000
	b4 bit = 0b0001_0000
	b5 bit = 0b0010_0000
	b6 bit = 0b0100_0000
	b7 bit = 0b1000_0000
)

func get(b byte, n bit) bool {
	return b&byte(n) > 0
}

func set(b byte, n bit) byte {
	return b | byte(n)
}

func clear(b byte, n bit) byte {
	return b ^ byte(n)
}
