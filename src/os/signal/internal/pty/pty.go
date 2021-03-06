// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris
// +build cgo

// Package pty is a simple pseudo-terminal package for Unix systems,
// implemented by calling C functions via cgo.
// This is only used for testing the os/signal package.
package pty

/*
#define _XOPEN_SOURCE 600
#include <fcntl.h>
#include <stdlib.h>
#include <unistd.h>
*/
import "C"

import (
	"fmt"
	"os"
)

// Open returns a master pty and the name of the linked slave tty.
func Open() (master *os.File, slave string, err error) {
	m, err := C.posix_openpt(C.O_RDWR)
	if err != nil {
		return nil, "", fmt.Errorf("posix_openpt: %v", err)
	}
	if _, err := C.grantpt(m); err != nil {
		C.close(m)
		return nil, "", fmt.Errorf("grantpt: %v", err)
	}
	if _, err := C.unlockpt(m); err != nil {
		C.close(m)
		return nil, "", fmt.Errorf("unlockpt: %v", err)
	}
	slave = C.GoString(C.ptsname(m))
	return os.NewFile(uintptr(m), "pty-master"), slave, nil
}
