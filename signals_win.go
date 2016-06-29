// +build windows

package main

import (
	"os"
	"syscall"
)

var signals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGQUIT,
	syscall.SIGTERM,
}
