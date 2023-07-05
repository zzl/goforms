package forms

import "log"

type SimpleWindowObject struct {
	WindowObject
	super *WindowObject
}

var _simpleWindowClassRegstered bool

func (this *SimpleWindowObject) EnsureClassRegistered() {
	if _simpleWindowClassRegstered {
		return
	}
	_, err := RegisterClass("goforms.simplewindow", nil, ClassOptions{
		BackgroundBrush: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	_simpleWindowClassRegstered = true
}

func (this *SimpleWindowObject) EnsureCustomWndProc() {
	//nop
}

func (this *SimpleWindowObject) GetWindowClass() string {
	return "goforms.simplewindow"
}
