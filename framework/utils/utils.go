package utils

import (
	"fmt"
	"github.com/zzl/goforms/framework/consts"
	"github.com/zzl/goforms/framework/types"
	"log"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Round rounds a float value to the nearest integer
// and returns it as the specified number type.
func Round[T types.Float, RT types.Number](value T) RT {
	return RT(math.Round(float64(value)))
}

// Abs returns the absolute value of a signed integer or float value.
func Abs[T types.SignedInt | types.Float](value T) T {
	if value < 0 {
		return -value
	}
	return value
}

// ToString converts a value to its string representation.
func ToString(value any) string {
	if IsNull(value) {
		return ""
	}
	s, ok := value.(string)
	if ok {
		return s
	}
	n, ok := value.(int)
	if ok {
		return strconv.Itoa(n)
	}
	return fmt.Sprintf("%v", value)
}

// ToTime converts a value to a Time object.
func ToTime(value any) time.Time {
	if value == nil {
		return time.Time{}
	}
	switch v := value.(type) {
	case string:
		return ParseTime(v)
	case int: //millis since Jan 1, 1970, 00:00:00, as in js
		return time.Unix(int64(v/1000), int64(v%1000*1000000))
	case time.Time:
		return v
	}
	return time.Time{}
}

// ParseTime parses a string and returns a Time object.
func ParseTime(str string) time.Time {
	if str == "" {
		return time.Time{}
	}
	if strings.Contains(str, "GMT") {
		tm, err := time.Parse("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)", str)
		if err != nil {
			return tm
		}
	}
	if matched, _ := regexp.MatchString(`\dT\d`, str); matched {
		tm, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return tm
		}
	}
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-1-2",
		"2006-1-2 15:04",
		"2006-1-2 15:04:05",
		"2006-1-2 15:4:5",
		"15:04",
		"15:04:05",
	}
	for _, format := range formats {
		tm, err := time.Parse(format, str)
		if err == nil {
			return tm
		}
	}
	return time.Time{}
}

// IsNull checks if a value is nil or equals to a special Null constant
func IsNull(value any) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case int32, uint32, int64, uint64, int, uint, float32, float64, types.Bool:
		return v == consts.Null
	case string:
		return v == consts.NullStr
	case nil:
		return true
	default:
		return reflect.ValueOf(value).IsNil()
	}
}

// AssignDefault updates a variable to the specified defaultValue
// if it's value is the zero value.
func AssignDefault[T types.Number | string](pValue *T, defaultValue T) {
	var value0 T
	if *pValue == value0 {
		*pValue = defaultValue
	}
}

// MagicZeroTo0 updates variables to 0 whose value is the special zero constant.
func MagicZeroTo0[T types.Numeric32BitsOrMore](pValues ...*T) {
	for _, pValue := range pValues {
		if *pValue == consts.Zero {
			*pValue = 0
		}
	}
}

// IfElse returns either ifValue or elseValue based on the condition.
func IfElse[T any](condition bool, ifValue T, elseValue T) T {
	if condition {
		return ifValue
	} else {
		return elseValue
	}
}

// If returns value if the condition is true,
// otherwise it returns the zero value of the type.
func If[T any](condition bool, value T) T {
	var result T
	if condition {
		result = value
	}
	return result
}

// OptionalArgByVal returns the first element of args if it exists,
// otherwise it returns the zero value of the type.
func OptionalArgByVal[T any](args []T) T {
	var t T
	if len(args) > 0 {
		t = args[0]
	}
	return t
}

// OptionalArg returns the first element of pointer args if it exists,
// otherwise it returns a new pointer to the zero value of the pointed type.
func OptionalArg[T any](args []*T) *T {
	var p *T
	if len(args) > 0 {
		p = args[0]
	} else {
		p = new(T)
	}
	return p
}

// AssertNoErr panics if an error is not nil.
func AssertNoErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
