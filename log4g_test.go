package log4g

import (
	"testing"
)

func TestLog4g(t *testing.T) {
	Init(Config{
		LogMode:   varMode,
		Path:      "logs",
		NameSpace: "knowing",
	})
	InfoFormat("info ")
	ErrorFormat("error ")
	Stat("stat ")
	Slow("Slow ")
}
