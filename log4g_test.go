package log4g

import (
	"testing"
)

func TestLog4g(t *testing.T) {
	InfoFormat("info ")
	ErrorFormat("error ")
	Stat("stat ")
	Slow("Slow ")
}
