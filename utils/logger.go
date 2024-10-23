package utils

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// terminal colors
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var faint = color.New(color.Faint).SprintFunc()

type log int

const (
	LOG_UNKNOWN log = iota
	LOG_INFO
	LOG_SUCCESS
	LOG_FAILURE
	LOG_STRICT_INFO
)

func Logger(s log, msg string) {
	statusToUnicode := map[log]string{
		LOG_INFO:        "",
		LOG_STRICT_INFO: green(">>"),
		LOG_SUCCESS:     green("\U00002714"),
		LOG_FAILURE:     red("\U00002718"),
	}

	if s == LOG_STRICT_INFO {
		msg = green(msg)
	} else if s != LOG_INFO && s != LOG_STRICT_INFO {
		msg = strings.ToLower(msg)             // convert to lowercase
		msg = strings.ReplaceAll(msg, ".", "") // remove punctuation
		msg = faint(msg)
	}

	fmt.Printf("%s %s\n", statusToUnicode[s], msg)
}

func ErrorSymbol() string {
	return red("\U00002718")
}
