package converter

import (
	"strconv"

	// Internal
	errDom "github.com/revirathya/revirathya-telbot/internal/domain/error"
)

// From int64 -> string
func FromInt64ToString(field string, v int64) string {
	s := strconv.FormatInt(v, 10)
	return s
}

func FromOptionalInt64ToString(field string, v *int64) *string {
	if v == nil {
		return nil
	}

	s := strconv.FormatInt(*v, 10)
	return &s
}

// From string -> int64
func FromStringToInt64(field string, v string) (int64, error) {
	num, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, &errDom.ValidationError{Field: field, Reason: "must be a valid integer"}
	}
	return num, nil
}

func FromOptionalStringToInt64(field string, v *string) (*int64, error) {
	if v == nil {
		return nil, nil
	}

	num, err := FromStringToInt64(field, *v)
	if err != nil {
		return nil, &errDom.ValidationError{Field: field, Reason: "must be a valid integer"}
	}
	return &num, nil
}

// From int -> string
func FromIntToString(field string, v int) string {
	return strconv.Itoa(v)
}

func FromOptionalIntToString(field string, v *int) *string {
	if v == nil {
		return nil
	}

	s := strconv.Itoa(*v)
	return &s
}

// From string -> int
func FromStringToInt(field string, v string) (int, error) {
	num, err := strconv.Atoi(v)
	if err != nil {
		return 0, &errDom.ValidationError{Field: field, Reason: "must be a valid integer"}
	}
	return num, nil
}

func FromOptionalStringToInt(field string, v *string) (*int, error) {
	if v == nil {
		return nil, nil
	}

	num, err := FromStringToInt(field, *v)
	if err != nil {
		return nil, &errDom.ValidationError{Field: field, Reason: "must be a valid integer"}
	}
	return &num, nil
}
