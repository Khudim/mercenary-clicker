package main

import (
	"bytes"
	"github.com/kbinani/screenshot"
	"github.com/lxn/win"
	"image/png"
	"io/ioutil"
	"syscall"
	"unsafe"
)

var (
	user32             = syscall.MustLoadDLL("user32.dll")
	procEnumWindows    = user32.MustFindProc("EnumWindows")
	procGetWindowTextW = user32.MustFindProc("GetWindowTextW")
)

type Window struct {
	x, y, width, height int
	h                   *syscall.Handle
	hwnd                win.HWND
}

func (w *Window) screenElement(t *Template) *Image {

	x := w.x + int(float32(w.width)*t.X)
	y := w.y + int(float32(w.height)*t.Y)
	width := int(float32(w.width) * t.Width)
	height := int(float32(w.height) * t.Height)

	img, _ := screenshot.Capture(x, y, width, height)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	if t.Debug {
		_ = ioutil.WriteFile("test-"+t.Name+".png", buf.Bytes(), 0644)
	}
	return &Image{x: x, y: y, width: width, height: height, bytes: buf.Bytes()}
}

func FindWindow(title string) (syscall.Handle, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		b := make([]uint16, 200)
		_, err := GetWindowText(h, &b[0], int32(len(b)))
		if err != nil {
			// ignore the error
			return 1 // continue enumeration
		}
		//fmt.Println(syscall.UTF16ToString(b))
		if syscall.UTF16ToString(b) == title {
			// note the window
			hwnd = h
			return 0 // stop enumeration
		}
		return 1 // continue enumeration
	})
	_ = EnumWindows(cb, 0)

	return hwnd, nil
}

func GetWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(procEnumWindows.Addr(), 2, enumFunc, lparam, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
