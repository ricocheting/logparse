package internal

// # go test -run=TestDumpRawBoltHitsData

// RAW dump the extensions in the database
import (
	"fmt"
	"testing"
)

func TestFormatShortHand(t *testing.T) {
	fmt.Println(FormatShortHand(12345678901))
	fmt.Println(FormatShortHand(2234567890))
	fmt.Println(FormatShortHand(323456789))
	fmt.Println(FormatShortHand(42345678))
	fmt.Println(FormatShortHand(5234567))
	fmt.Println(FormatShortHand(623456))
	fmt.Println(FormatShortHand(72345))
	fmt.Println(FormatShortHand(8234))
	fmt.Println(FormatShortHand(923))
	fmt.Println(FormatShortHand(99999999))
}

func TestFormatMonth(t *testing.T) {

	fmt.Println(FormatMonth(0))
	fmt.Println(FormatMonth(1))
	fmt.Println(FormatMonth(2))
	fmt.Println(FormatMonth(3))
	fmt.Println(FormatMonth(4))
	fmt.Println(FormatMonth(5))
	fmt.Println(FormatMonth(6))
	fmt.Println(FormatMonth(7))
	fmt.Println(FormatMonth(8))
	fmt.Println(FormatMonth(9))
	fmt.Println(FormatMonth(10))
	fmt.Println(FormatMonth(11))
	fmt.Println(FormatMonth(12))
	fmt.Println(FormatMonth(13))
}
