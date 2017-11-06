package internal

import "encoding/binary"

func (s Stats) ToSlice(min uint64) []Stat {
	out := make([]Stat, 0, len(s))
	for k, v := range s {
		if min > 0 && v < min {
			continue
		}
		out = append(out, Stat{k, v})
	}
	return out[:len(out):len(out)] // trim the slice to release the unused memory
}

// itob converts an uint64 to byte array
func Itob(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// btoi converts a byte array to uint64
func Btoi(b []byte) uint64 {
	i := binary.BigEndian.Uint64(b)
	return i
}
