package main

// #include "Windows.h"
import "C"
import (
	"fmt"
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

type callback func(uint64, uint64)

type Listener struct {
	functions []callback
	hook      uintptr
}

const (
	WM_KEYDOWN = C.WM_KEYDOWN
	WM_KEYUP   = C.WM_KEYUP
)

var (
	user32              = windows.NewLazyDLL("user32.dll")
	setWindowsHookEx    = user32.NewProc("SetWindowsHookExA")
	unhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	getMessage          = user32.NewProc("GetMessageA")
	translateMessage    = user32.NewProc("TranslateMessage")
	dispatchMessage     = user32.NewProc("DispatchMessage")
)

func NewListener() (*Listener, error) {
	listener := new(Listener)

	callback := windows.NewCallback(func(_, wParam, lParam uintptr) uintptr {
		t := (*C.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		for _, f := range listener.functions {
			f(uint64(wParam), uint64(t.vkCode))
		}

		return 0
	})

	hHook, code, err := setWindowsHookEx.Call(C.WH_KEYBOARD_LL, callback, 0, 0)

	if code != 0 {
		return nil, err
	}

	listener.hook = hHook

	return listener, nil
}

func (l *Listener) Add(f callback) {
	l.functions = append(l.functions, f)
}

func (l *Listener) Listen() {
	var (
		msg    C.MSG
		msgPtr = uintptr(unsafe.Pointer(&msg))
		nilPtr = uintptr(unsafe.Pointer(nil))
	)

	r, _, err := getMessage.Call(msgPtr, nilPtr, nilPtr, nilPtr)

	for ; r != 0; r, _, err = getMessage.Call(msgPtr, nilPtr, nilPtr, nilPtr) {
		if (int)(r) == -1 {
			log.Print(err)
			break
		} else {
			translateMessage.Call(msgPtr)
			dispatchMessage.Call(msgPtr)
		}
	}

	unhookWindowsHookEx.Call(l.hook)
}
