package main

// #include "Windows.h"
import "C"
import (
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

type callback func(uintptr, C.ULONG)

type Listener struct {
	functions []callback
	hook      uintptr
}

var (
	user32              = windows.NewLazyDLL("user32.dll")
	SetWindowsHookEx    = user32.NewProc("SetWindowsHookExA")
	UnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	GetMessage          = user32.NewProc("GetMessageA")
	TranslateMessage    = user32.NewProc("TranslateMessage")
	DispatchMessage     = user32.NewProc("DispatchMessage")
)

func NewListener() (*Listener, error) {
	listener := new(Listener)

	callback := windows.NewCallback(func(_, wParam, lParam uintptr) uintptr {
		t := (*C.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		for _, f := range listener.functions {
			f(wParam, t.vkCode)
		}

		return 0
	})

	hHook, code, err := SetWindowsHookEx.Call(C.WH_KEYBOARD_LL, callback, 0, 0)

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

	r, _, err := GetMessage.Call(msgPtr, nilPtr, nilPtr, nilPtr)

	for ; r != 0; r, _, err = GetMessage.Call(msgPtr, nilPtr, nilPtr, nilPtr) {
		if (int)(r) == -1 {
			log.Print(err)
			break
		} else {
			TranslateMessage.Call(msgPtr)
			DispatchMessage.Call(msgPtr)
		}
	}

	UnhookWindowsHookEx.Call(l.hook)
}
