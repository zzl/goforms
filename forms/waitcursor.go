package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
)

type WaitCursor struct {
	oriWaitCursor *WaitCursor
	hCursor       win32.HCURSOR
	hOriCursor    win32.HCURSOR
}

var curWaitCursor *WaitCursor

func NewWaitCursor() WaitCursor {
	hCursor, _ := win32.LoadCursor(0, win32.IDC_WAIT)
	hOriCursor := win32.SetCursor(hCursor)
	wc := WaitCursor{
		hCursor:       hCursor,
		hOriCursor:    hOriCursor,
		oriWaitCursor: curWaitCursor,
	}
	curWaitCursor = &wc
	return wc
}

func (me WaitCursor) Update() {
	win32.SetCursor(me.hCursor)
}

func (me WaitCursor) Restore() {
	Dispatcher.Invoke(func() {
		curWaitCursor = me.oriWaitCursor
		win32.SetCursor(me.hOriCursor)
	})
}
