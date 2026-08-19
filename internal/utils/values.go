package utils

const EMPTY_CALLBACK string = "em"

func EmptyCallback() *string {
	val := EMPTY_CALLBACK
	return &val
}
