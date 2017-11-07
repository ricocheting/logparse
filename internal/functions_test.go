package internal

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
