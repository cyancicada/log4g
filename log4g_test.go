package log4g

import (
	"fmt"
	"testing"
)

func TestLog4g(t *testing.T) {
	Init(Config{Path: "logs", Stdout: true})
	Slow("Slow ")
	fmt.Println(Log(1, 2, "a", "b"))
}
