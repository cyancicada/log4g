package log4g

import (
	"testing"
)

func TestLog4g(t *testing.T) {
	Init(Config{Path: "logs", Stdout: true})
	Slow("Slow ")
	Stat("a")
}
