package utils

import (
	"fmt"

	"github.com/fatih/color"
)

// terminal colors
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

type status int

const (
	STATUS_UNKNOWN status = iota
	STATUS_INFO
	STATUS_SUCCESS
	STATUS_FAILURE
)

func Logger(s status, log string) {
	statusToUnicode := map[status]string{
		STATUS_INFO:    "",
		STATUS_SUCCESS: green("\U00002714 "),
		STATUS_FAILURE: red("\U00002717 "),
	}

	fmt.Printf("%s%s\n", statusToUnicode[s], log)
}
