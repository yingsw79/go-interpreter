package util

import (
	"fmt"
	"strings"
)

var traceLevel int

const traceIdentPlaceholder string = "\t"

func identLevel() string {
	return strings.Repeat(traceIdentPlaceholder, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", identLevel(), fs)
}

func Trace(msg string) string {
	traceLevel++
	tracePrint("BEGIN " + msg)
	return msg
}

func Untrace(msg string) {
	tracePrint("END " + msg)
	traceLevel--
}
