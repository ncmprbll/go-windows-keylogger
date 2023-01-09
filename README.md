# Windows Keylogger written in Go
Simple keylogging library using Windows hooks.

# Why?
I have found other Windows solutions to be ineffective and resource-demanding (especially regarding CPU) since they were using while loops, loops inside of while loops and some other forbidden wizardry.

# Example
```Go
package main

import (
	"fmt"

	"github.com/ncmprbll/go-windows-keylogger"
)

func main() {
	listener, err := keylogger.NewListener()

	if err != nil {
		return
	}

	// vkCode == 112 is F1

	listener.Add(func(wParam uint64, vkCode uint64) {
		if wParam == keylogger.WM_KEYDOWN && vkCode == 112 {
			fmt.Println("F1 is being pressed or held down!")
		}
	})

	listener.Add(func(wParam uint64, vkCode uint64) {
		if wParam == keylogger.WM_KEYUP && vkCode == 112 {
			fmt.Println("Thank you for releasing the F1 button!")
		}
	})

	listener.Listen()
}
```

<p align="center">
  <img src="https://i.imgur.com/jXZievS.gif" alt="animated"/>
</p>
