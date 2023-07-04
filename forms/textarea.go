package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type TextArea interface {
	Edit

	TextAreaObj() *TextAreaObject
}

type TextAreaObject struct {
	EditObject
	super *EditObject

	NoAutoWrap          bool //horz scrollbar?
	AlwaysShowScrollBar bool
	ForeColorAwareSupport

	paddingTop    int
	paddingBottom int
	lineHeight    int
}

type NewTextArea struct {
	Parent Container
	Name   string
	Text   string
	Pos    Point
	Size   Size
}

func (me NewTextArea) Create(extraOpts ...*WindowOptions) TextArea {
	textArea := NewTextAreaObject()
	textArea.name = me.Name

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Text
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y

	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := textArea.Create(*opts)
	assertNoErr(err)
	configControlSize(textArea, me.Size)

	return textArea
}

func NewTextAreaObject() *TextAreaObject {
	return virtual.New[TextAreaObject]()
}

func textAreaWndProc(win *WindowObject, m *Message) error {
	if m.UMsg == win32.WM_GETDLGCODE { //??
		return m.SetHandledWithResult(win32.LRESULT(win32.DLGC_WANTALLKEYS))
	} else if m.UMsg == win32.WM_CHAR {
		if m.WParam == '\t' {
			HandleTabFocus(m.HWnd)
			return m.SetHandledWithResult(0)
		}
	}
	//???
	return nil
}

func (this *TextAreaObject) TextAreaObj() *TextAreaObject {
	return this
}

func (this *TextAreaObject) Init() {
	this.super.Init()
	this.WinProcFunc = textAreaWndProc
}

func (this *TextAreaObject) Dispose() {
	this.super.Dispose()
}

func (this *TextAreaObject) GetWindowClass() string {
	return "Edit"
}

func (this *TextAreaObject) GetPreferredSize(maxCx int, maxCy int) (int, int) {
	_, h := MeasureText(this.Handle, "|Why")
	return 32, (h+2)*2 + 2
}

func (this *TextAreaObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	include = WINDOW_STYLE(win32.ES_MULTILINE | win32.ES_AUTOVSCROLL | win32.ES_WANTRETURN)
	if this.NoAutoWrap {
		include |= WINDOW_STYLE(win32.ES_AUTOHSCROLL)
	}
	return
}

func (this *TextAreaObject) Create(options WindowOptions) error {
	err := this.super.Create(options)
	if err != nil {
		return err
	}

	//
	var rc win32.RECT
	win32.GetWindowRect(this.Handle, &rc)

	win32.SetWindowPos(this.Handle, 0, 0, 0, 40, 40,
		win32.SWP_NOREDRAW|win32.SWP_NOZORDER|win32.SWP_NOMOVE)

	var rcText win32.RECT
	SendMessage(this.Handle, win32.EM_GETRECT, 0, unsafe.Pointer(&rcText))

	var rcCli win32.RECT
	win32.GetClientRect(this.Handle, &rcCli)

	this.paddingTop = int(rcText.Top)
	this.paddingBottom = int(rcCli.Bottom - rcText.Bottom)

	SetWindowText(this.Handle, "|\r\n|")
	ret, _ := SendMessage(this.Handle, win32.EM_POSFROMCHAR, 0, 0)
	_, y1 := win32.LOWORD(win32.DWORD(ret)), int16(win32.HIWORD(win32.DWORD(ret)))
	ret, _ = SendMessage(this.Handle, win32.EM_POSFROMCHAR, 3, 0)
	_, y2 := win32.LOWORD(win32.DWORD(ret)), int16(win32.HIWORD(win32.DWORD(ret)))
	SetWindowText(this.Handle, "")

	this.lineHeight = int(y2 - y1)
	//println(lineHeight)

	win32.SetWindowPos(this.Handle, 0, 0, 0,
		rc.Right-rc.Left, rc.Bottom-rc.Top,
		win32.SWP_NOREDRAW|win32.SWP_NOZORDER|win32.SWP_NOMOVE)

	//
	return err
}

func (this *TextAreaObject) SetBounds(left, top, width, height int) {
	this.super.SetBounds(left, top, width, height)
	this.updateScrollbar()
}

func (this *TextAreaObject) updateScrollbar() {
	if this.AlwaysShowScrollBar {
		win32.ShowScrollBar(this.Handle, win32.SB_VERT, 1)
		return
	}

	var rcCli win32.RECT
	win32.GetClientRect(this.Handle, &rcCli)

	lineCount, _ := SendMessage(this.Handle, win32.EM_GETLINECOUNT, 0, 0)

	textHeight := this.lineHeight*int(lineCount) + this.paddingTop + this.paddingBottom
	if textHeight > int(rcCli.Bottom) {
		win32.ShowScrollBar(this.Handle, win32.SB_VERT, 1)
	} else {
		win32.ShowScrollBar(this.Handle, win32.SB_VERT, 0)
	}
}

func (this *TextAreaObject) OnReflectCommand(msg *CommandMessage) {
	if msg.GetNotifyCode() == uint16(win32.EN_CHANGE) {
		this.updateScrollbar()
	}
}
