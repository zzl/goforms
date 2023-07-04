package forms

import (
	"fmt"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/utils"
	"log"
)

type dispatcherImpl struct {
	threadId uint32
	chAction chan Action
}

func newDispatcherImpl() *dispatcherImpl {
	d := &dispatcherImpl{}
	d.chAction = make(chan Action, 256)
	d.threadId = win32.GetCurrentThreadId()
	return d
}

// Dispatcher is used to dispatch actions to be executed on the UI thread
var Dispatcher = newDispatcherImpl()

// Invoke executes the action on the UI thread
func (this *dispatcherImpl) Invoke(action Action, optSync ...bool) {
	syncWait := utils.OptionalArgByVal(optSync)
	var chWait chan struct{}
	if syncWait {
		if this.threadId == win32.GetCurrentThreadId() {
			action()
			return //
		}
		action0 := action
		chWait = make(chan struct{})
		action = func() {
			action0()
			close(chWait)
		}
	}
	this.chAction <- action
	ok, errno := win32.PostThreadMessage(this.threadId, WM_APP_DISPATCH, 0, 0)
	if ok != win32.TRUE {
		log.Panic(errno)
	}
	if chWait != nil {
		<-chWait
	}
}

func (this *dispatcherImpl) check() {
	for {
		select {
		case action := <-this.chAction:
			action()
		default:
			return
		}
	}
}

// MsgPreprocessor is an interface that wraps the PreprocessMsg method.
type MsgPreprocessor interface {
	// PreprocessMsg is called to preprocess win32 messages in the message loop
	PreprocessMsg(msg *win32.MSG) bool
}

type msgPreprocessorsImpl struct {
	processors []MsgPreprocessor
}

func (this *msgPreprocessorsImpl) Add(processor MsgPreprocessor) {
	this.processors = append(this.processors, processor)
}

func (this *msgPreprocessorsImpl) Remove(processor MsgPreprocessor) {
	count := len(this.processors)
	for n := 0; n < count; n++ {
		if this.processors[n] == processor {
			this.processors = append(this.processors[:n], this.processors[n+1:]...)
			count -= 1
			break
		}
	}
	if count == 0 {
		this.processors = nil
	}
}

func (this *msgPreprocessorsImpl) PreprocessMsg(msg *win32.MSG) bool {
	for _, processor := range this.processors {
		if processor.PreprocessMsg(msg) {
			return true
		}
	}
	return false
}

// MsgPreprocessors is a global registry to add or remove MsgPreprocessor implementations
var MsgPreprocessors msgPreprocessorsImpl

// MessageLoop is the function that's typically called at the end of the main function
// to pump windows messages to their desired target and bring the UI to life.
func MessageLoop() {
	var msg win32.MSG
	for {
		bRet, _ := win32.GetMessage(&msg, 0, 0, 0)
		if bRet == 0 { //WM_QUIT
			break
		}
		if bRet == -1 {
			fmt.Println("??")
			break
		}
		processMsg(&msg)
	}
}

func processMsg(msg *win32.MSG) {
	if msg.Message == WM_APP_DISPATCH {
		Dispatcher.check()
		return
	}
	toDispatch := true
	if HWndActive != 0 {
		if hAccelActive != 0 {
			translated, _ := win32.TranslateAccelerator(
				HWndActive, hAccelActive, msg)
			if translated != 0 {
				toDispatch = false
			}
		}
	}
	if toDispatch {
		if MsgPreprocessors.PreprocessMsg(msg) {
			toDispatch = false
		}
	}
	if toDispatch {
		win, ok := windowMap[msg.Hwnd]
		if ok && win.PreProcessMsg(msg) {
			toDispatch = false
		}
	}
	if toDispatch {
		if msg.Message == win32.WM_SYSCHAR && FindFirstFocusable(HWndActive) == 0 {
			//nop
		} else {
			isDlgMsg := win32.IsDialogMessage(HWndActive, msg)
			if isDlgMsg != win32.FALSE {
				toDispatch = false
			}
		}
	}

	if toDispatch {
		win32.TranslateMessage(msg)
		win32.DispatchMessage(msg)
	}
}

// DoEvents processes all Windows messages currently in the message queue.
func DoEvents() {
	var msg win32.MSG
	for {
		ret := win32.PeekMessage(&msg, 0, 0, 0xFFFF, win32.PM_REMOVE)
		if ret == 0 { //mq empty
			break
		} else {
			processMsg(&msg)
		}
	}
}
