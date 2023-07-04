package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/events"
	"log"
	"syscall"
)

type Timer struct {
	IntervalMillis int
	OnTick         SimpleEvent

	id uintptr
}

func NewTimer() *Timer {
	return &Timer{}
}

var _timerMap = map[uintptr]*Timer{}

var _pTimerProc uintptr

func _timerProc(hWnd win32.HWND, uMsg uint32, idEvent uintptr, dwTime win32.DWORD) uintptr {
	timer := _timerMap[idEvent]
	timer.OnTick.Fire(timer, &SimpleEventInfo{})
	return 0
}

func (this *Timer) Start() {
	if _pTimerProc == 0 {
		_pTimerProc = syscall.NewCallback(_timerProc)
	}
	id, errno := win32.SetTimer(0, 0, uint32(this.IntervalMillis), _pTimerProc)
	if id == 0 {
		log.Panic(errno)
	}
	this.id = id
	_timerMap[id] = this
}

func (this *Timer) Dispose() {
	win32.KillTimer(0, this.id)
	delete(_timerMap, this.id)
}

func CreateTimer(intervalMillis int, onTick events.SimpleEventListener) *Timer {
	timer := NewTimer()
	timer.IntervalMillis = intervalMillis
	timer.OnTick.AddListener(onTick)
	return timer
}
