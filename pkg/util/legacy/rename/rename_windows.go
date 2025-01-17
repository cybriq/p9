package rename

// The following is adapted from goleveldb
// (https://github.com/syndtr/goleveldb) under the following license:
//
// Copyright 2012 Suryandaru Triandana <syndtr@gmail.com> All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the
// following conditions are met:
//
//     * Redistributions of source code must retain the above copyright notice, this list of conditions and the following
//     disclaimer.
//
//     * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following
//     disclaimer in the documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.import (
import (
	"syscall"
	"unsafe"
)

var (
	modkernel32     = syscall.NewLazyDLL("kernel32.dll")
	procMoveFileExW = modkernel32.NewProc("MoveFileExW")
)

const (
	_MOVEFILE_REPLACE_EXISTING = 1
)

func moveFileEx(from *uint16, to *uint16, flags uint32) (e error) {
	r1, _, e1 := syscall.Syscall(
		procMoveFileExW.Addr(), 3,
		uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(to)),
		uintptr(flags),
	)
	if r1 == 0 {
		if e1 != 0 {
			return error(e1)
		} else {
			return syscall.EINVAL
		}
	}
	return nil
}

// Atomic provides an atomic file rename. newpath is replaced if it already exists.
func Atomic(oldpath, newpath string) (e error) {
	from, e := syscall.UTF16PtrFromString(oldpath)
	if e != nil {
		E.Ln(e)
		return e
	}
	to, e := syscall.UTF16PtrFromString(newpath)
	if e != nil {
		E.Ln(e)
		return e
	}
	return moveFileEx(from, to, _MOVEFILE_REPLACE_EXISTING)
}
