package log4g

import (
	"testing"
)

func TestLog4g(t *testing.T) {
	Init(Config{
		Path:   "logs",
		Stdout: true,
	})
	InfoFormat("info ")
	ErrorFormat("error ")
	Stat("stat ")
	Slow("Slow ")
}
