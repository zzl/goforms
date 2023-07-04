package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"syscall"
	"time"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type DateTimePicker interface {
	Control
	Input
}

type DateTimePickerObject struct {
	ControlObject
	super *ControlObject

	Nullable  bool
	Format    string
	TimeOnly  bool
	UseUpDown bool

	onValueChange SimpleEvent
}

func (this *DateTimePickerObject) GetValue() any {
	var st win32.SYSTEMTIME
	ret, errno := SendMessage(this.Handle, win32.DTM_GETSYSTEMTIME,
		0, unsafe.Pointer(&st))
	if int32(ret) == win32.GDT_ERROR {
		log.Fatal(errno)
	}
	if win32.NMDATETIMECHANGE_FLAGS(ret) == win32.GDT_NONE {
		return nil
	}
	tm := time.Date(int(st.WYear), time.Month(st.WMonth), int(st.WDay),
		int(st.WHour), int(st.WMinute), int(st.WSecond), 0, time.Local)
	return tm
	//return
}

func (this *DateTimePickerObject) SetValue(value any) {
	tm := utils.ToTime(value)
	if value == nil {
		if !this.Nullable {
			//?
			return
		}
		SendMessage(this.Handle, win32.DTM_SETSYSTEMTIME, (uint32)(win32.GDT_NONE), 0)
		return
	}
	var st win32.SYSTEMTIME
	st.WYear = uint16(tm.Year())
	st.WMonth = uint16(tm.Month())
	st.WDay = uint16(tm.Day())
	st.WHour = uint16(tm.Hour())
	st.WMinute = uint16(tm.Minute())
	st.WSecond = uint16(tm.Second())
	SendMessage(this.Handle, win32.DTM_SETSYSTEMTIME,
		(uint32)(win32.GDT_VALID), unsafe.Pointer(&st))
}

func (this *DateTimePickerObject) GetOnValueChange() *SimpleEvent {
	return &this.onValueChange
}

func (this *DateTimePickerObject) GetWindowClass() string {
	return "SysDateTimePick32"
}

func NewDateTimePickerObject() *DateTimePickerObject {
	return virtual.New[DateTimePickerObject]()
}

func (this *DateTimePickerObject) Init() {
	this.super.Init()
}

func (this *DateTimePickerObject) OnReflectNotify(msg *NotifyMessage) {
	if msg.GetNMHDR().Code == uint32(win32.DTN_DATETIMECHANGE) {
		this.onValueChange.Fire(this, &SimpleEventInfo{})
		msg.Handled = true
		msg.Result = 0
	}
}

func (this *DateTimePickerObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	var style WINDOW_STYLE
	if this.Nullable {
		style |= WINDOW_STYLE(win32.DTS_SHOWNONE)
	}
	if this.UseUpDown {
		style |= WINDOW_STYLE(win32.DTS_UPDOWN)
	}
	if this.TimeOnly {
		style |= WINDOW_STYLE(win32.DTS_TIMEFORMAT)
	}
	return style, 0
}

func (this *DateTimePickerObject) Create(options WindowOptions) error {
	err := this.super.Create(options)
	if this.Format != "" {
		pwstr, _ := syscall.UTF16PtrFromString(this.Format)
		SendMessage(this.Handle, win32.DTM_SETFORMAT,
			0, unsafe.Pointer(pwstr))
	}
	return err
}

func (this *DateTimePickerObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	var sz win32.SIZE
	bOk, errno := SendMessage(this.Handle, win32.DTM_GETIDEALSIZE,
		0, unsafe.Pointer(&sz))
	if bOk == 0 {
		log.Println(errno)
	}
	_, cyText := MeasureText(this.Handle, "|Why")
	cyBorder, _ := win32.GetSystemMetrics(win32.SM_CYBORDER)
	//from winform source
	height := cyText + int(cyBorder)*4 + 3
	return int(sz.Cx), height
}
