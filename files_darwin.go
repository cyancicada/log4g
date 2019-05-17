package log4g

import (
	"os"
	"syscall"
)

func CloseOnExec(file *os.File) {
	if file != nil {
		syscall.CloseOnExec(int(file.Fd()))
	}
}
