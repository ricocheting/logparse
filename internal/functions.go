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

//IsNewerDay checks (while ignoring timezone) has the date rolled over
func IsNewerDay(startDate, compDate time.Time) bool {
	/*y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2*/

	d1 := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(compDate.Year(), compDate.Month(), compDate.Day(), 0, 0, 0, 0, time.UTC)

	return !(d1.Equal(d2) || d2.Before(d1))
}

// Itob converts an uint64 to byte array
func Itob(i uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return b
}

// Btoi converts a byte array to uint64. that's not a string number, that's the binary def of the number. just make sure that it doesn't bite you in the ass
func Btoi64(b []byte) uint64 {
	i := binary.BigEndian.Uint64(b)
	return i
}
func Btoi8(b []byte) uint8 {
	n, _ := strconv.Atoi(string(b))
	return uint8(n)
}
func Btoi16(b []byte) uint16 {
	n, _ := strconv.Atoi(string(b))
	return uint16(n)
}
func Atoi8(s string) uint8 {
	n, _ := strconv.Atoi(s)
	return uint8(n)
}

// FormatShortDate turns "YYYYMMDD" into "Jan 01"
func FormatShortDate(ts string) string {
	out, _ := time.Parse("20060102", ts)
	return out.Format("Jan 02")
}

// FormatMonth turns "MM" into "January"
func FormatMonth(m uint8) string {
	out, _ := time.Parse("01", strconv.Itoa(int(m)))
	return out.Format("January")
}

// FormatShortDate turns "D" into "Jan 01"
func FormatShortMonth(ts time.Time) string {
	//out, _ := time.Parse("1", strconv.Itoa(int(m)))
	return ts.Format("Jan")
}

// FormatMonthYear turns timestamp into "January 2017"
func FormatMonthYear(ts time.Time) string {
	//out, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", ts)
	return ts.Format("January 2006")
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

// FormatStatusCodeName turns "404" into "Not Found"
func FormatStatusCodeName(code string) string {
	if val, ok := StatusCodeNames[code]; ok {
		return val
	} else {
		return ""
	}
}

// PathDirectory turns YYYY into "2017"
func PathDirectory(y uint16) string {
	return strconv.Itoa(int(y))
}

// PathFilename turns MM into "01-January.html"
func PathFilename(m uint8) string {
	return strconv.Itoa(int(m)) + "-" + FormatMonth(m) + ".html"
}
