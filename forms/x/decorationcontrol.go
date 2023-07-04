package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	. "github.com/zzl/goforms/forms"
)

type DecorationControl interface {
	Control
}

type DecorationControlObject struct {
	ControlObject
}

func (this *DecorationControlObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	return 0, win32.WS_TABSTOP
}
