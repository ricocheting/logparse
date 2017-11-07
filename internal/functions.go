package internal

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"
)

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

// Itob converts an uint64 to byte array
func Itob(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// Btoi converts a byte array to uint64
func Btoi(b []byte) uint64 {
	i := binary.BigEndian.Uint64(b)
	return i
}

// FormatShortDate turns "YYYYMMDD" into "Jan 01"
func FormatShortDate(ts string) string {
	out, _ := time.Parse("20060102", ts)
	return out.Format("Jan 01")
}

// FormatCommas turns 1234567890 into 1,234,567,890
func FormatCommas(n uint64) string {
	in := strconv.FormatUint(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

// FormatShortHand rounds and turns 42345678 into more human readable 42.3m
func FormatShortHand(in uint64) string {

	sz := float64(in)
	switch {
	case (in > 1000000000): //>1b = 1.23b
		return fmt.Sprintf("%.2fb", (sz / 1000000000))
	case (in > 100000000): //>100m 123m
		return fmt.Sprintf("%.0fm", (sz / 1000000))
	case (in > 10000000): //>10m 12.3m
		return fmt.Sprintf("%.1fm", (sz / 1000000))
	case (in > 1000000): //>1m 1.23m
		return fmt.Sprintf("%.2fm", (sz / 1000000))
	case (in > 100000): //>100k 123k
		return fmt.Sprintf("%.0fk", (sz / 1000))
	case (in > 10000): //>10k 12.3k
		return fmt.Sprintf("%.1fk", (sz / 1000))
	case (in > 1000): //>1k 1.23k
		return fmt.Sprintf("%.2fk", (sz / 1000))
	}

	return fmt.Sprintf("%d", in)
}