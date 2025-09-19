package file

import "github.com/DesSolo/rtc/pkg/rtc"

// Reader ...
type Reader interface {
	Read([]byte) (map[rtc.Key]rtc.Value, error)
}

// ReaderFunc ...
type ReaderFunc func([]byte) (map[rtc.Key]rtc.Value, error)

// Read ...
func (f ReaderFunc) Read(data []byte) (map[rtc.Key]rtc.Value, error) {
	return f(data)
}
