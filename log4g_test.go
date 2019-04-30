package log4g

import (
	"testing"
)

func TestLog4g(t *testing.T) {
	Init(Config{Path: "log", Stdout: true})
	Slow("Slow ")
}
