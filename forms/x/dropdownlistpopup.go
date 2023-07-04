package forms

import "github.com/zzl/go-win32api/v2/win32"
import . "github.com/zzl/goforms/forms"

type DropdownListPopup struct {
	ListBox ListBox

	onOk     SimpleEvent
	onCancel SimpleEvent
}

func (this *DropdownListPopup) ProcessMessage(m *Message) {
	//println("$", msg.Message)
	if m.UMsg == win32.WM_KEYDOWN {
		key := win32.VIRTUAL_KEY(m.WParam)
		if key == win32.VK_RETURN {
			this.onOk.Fire(this, nil)
			m.SetHandled(true)
		} else if key == win32.VK_ESCAPE || key == win32.VK_F4 {
			this.onCancel.Fire(this, nil)
			m.SetHandled(true)
		}
	} else if m.UMsg == win32.WM_LBUTTONUP {
		x, y, _ := ParseMouseMsgParams(m.WParam, m.LParam)
		//hWndLb := this.ListBox.GetHandle()
		//lparam := win32.LPARAM(win32.MAKELONG(uint16(x), uint16(y)))
		//ret, _ := win32.SendMessage(hWndLb, win32.LB_ITEMFROMPOINT, 0, lparam)
		//dwRet := win32.DWORD(ret)
		//index := win32.LOWORD(dwRet)
		//_= index
		//outside := win32.HIWORD(dwRet)
		//if outside == 0 {
		//	this.onOk.Fire(this, nil)
		//}
		index := this.ListBox.IndexFromPoint(int(x), int(y))
		if index != -1 {
			this.onOk.Fire(this, nil)
		}
	}
}

func (this *DropdownListPopup) Init() {
	lb := this.ListBox
	//lb.GetEvent(win32.WM_LBUTTONUP).AddListener(func(ei *EventInfo) {
	//	info := ToMsgEventInfo(ei)
	//	x, y, _ := ParseMouseMsgParams(info.WParam, info.LParam)
	//	var rc win32.RECT
	//	win32.GetClientRect(this.ListBox.GetHandle(), &rc)
	//	if PtInRect(&rc, win32.POINT{int32(x), int32(y)}) {
	//		this.onOk.Fire(this, nil)
	//	}
	//})
	//lb.GetEvent(win32.WM_KEYDOWN).AddListener(func(ei *EventInfo) {
	//	info := ToKeyEventInfo(ei)
	//	key := info.GetKey()
	//	if key == win32.VK_RETURN {
	//		this.onOk.Fire(this, nil)
	//		ei.SetHandled()
	//	} else if key == win32.VK_ESCAPE || key == win32.VK_F4 {
	//		this.onCancel.Fire(this, nil)
	//		ei.SetHandled()
	//	}
	//})
	//lb.AddMsgFilter(this)
	lb.AddMessageProcessor(this)
}

func (this *DropdownListPopup) Dispose() {
	if this.ListBox == nil {
		return
	}
	this.ListBox.RemoveMessageProcessor(this)
	this.ListBox = nil
}

func NewDropdownListPopup(listBox ListBox) *DropdownListPopup {
	obj := &DropdownListPopup{ListBox: listBox}
	obj.Init()
	return obj
}

func (this *DropdownListPopup) GetOnOk() *SimpleEvent {
	return &this.onOk
}

func (this *DropdownListPopup) GetOnCancel() *SimpleEvent {
	return &this.onCancel
}

func (this *DropdownListPopup) GetValue() interface{} {
	return this.ListBox.GetValue()
}

func (this *DropdownListPopup) GetText() string {
	return this.ListBox.GetText()
}

func (this *DropdownListPopup) GetControl() Control {
	return this.ListBox
}

func (this *DropdownListPopup) SetValue(value interface{}) {
	this.ListBox.SetValue(value)
}

func (this *DropdownListPopup) PreparePopup() {
	//
}

func (this *DropdownListPopup) GetPopupSize(width int,
	maxWidth int, maxHeight int) (int, int) {

	cx, cy := this.ListBox.GetPreferredSize(4096, 4096)
	if cx < width {
		cx = width
	}
	var cyNc int
	dwExStyle, _ := win32.GetWindowLong(this.ListBox.GetHandle(), win32.GWL_EXSTYLE)
	if (dwExStyle & int32(win32.WS_EX_CLIENTEDGE)) != 0 {
		cyNc = 4
	}
	return cx, cy + cyNc
}

func (this *DropdownListPopup) NotifyBeforeShow() {
	//
}

func (this *DropdownListPopup) NotifyAfterShow() {
	//
}

func (this *DropdownListPopup) SetContainer(container DropdownPopupContainer) {
	//
}
