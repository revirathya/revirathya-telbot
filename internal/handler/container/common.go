package container

import "time"

type InputDate struct {
	IsFill bool
	Value  time.Time
}

type InputString struct {
	IsFill bool
	Value  string
}

type InputInt struct {
	IsFill bool
	Value  int64
}
