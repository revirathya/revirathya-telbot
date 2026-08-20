package utils

import "github.com/dustin/go-humanize"

func FormatThousand(value int64) string {
	return humanize.Comma(value)
}
