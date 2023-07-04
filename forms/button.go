package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"unsafe"

	"github.com/zzl/goforms/framework/virtual"

	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing/colors"
)

// Button is an interface that represents a button control.
// It extends the Control interfaces.
type Button interface {
	Control // the parent interface

	ActionAware    // action aware
	CommandAware   // command aware
	TextAware      // text aware
	ForeColorAware // forecolor aware

	SetDefault()              //sets the button as the default button
	SetIcon(icon win32.HICON) // sets the icon for the button
	Click()                   // simulates a button click

	ButtonObj() *ButtonObject // returns the underlying ButtonObject
}

// ButtonSpi is an interface that provides additional methods
// specific to implementing a ButtonS.
type ButtonSpi interface {
	ControlSpi
	OnClick() // todo:
}

// ButtonInterface is a composition of Button and ButtonSpi
type ButtonInterface interface {
	Button
	ButtonSpi
}

// ButtonObject implements the ButtonInterface
// It extends ControlObject.
type ButtonObject struct {

	// ControlObject is the parent struct.
	ControlObject

	// super is the special pointer to the parent struct.
	super *ControlObject

	// ActionAwareSupport is the ActionAware implementation
	ActionAwareSupport

	// ForeColorAwareSupport is the ForeColorAware implementation
	ForeColorAwareSupport

	// command is the Command associated with the button
	command *Command
}

// NewButtonObject creates a new ButtonObject.
func NewButtonObject() *ButtonObject {
	return virtual.New[ButtonObject]()
}

// NewButton is a struct representing the configuration for creating a new button
type NewButton struct {
	Parent   Container   // Parent container
	Id       uint16      // Control identifier
	Name     string      // Control Name
	Text     string      // Text for the button
	Icon     win32.HICON // Icon of the button
	Pos      Point       // Control position
	Size     Size        // Control Size
	Disabled bool        // Whether the button is disabled
	Action   Action      // Action to execute on click
}

// Create creates a new button with the specified configuration
func (me NewButton) Create(extraOpts ...*WindowOptions) Button {
	btn := NewButtonObject()
	btn.name = me.Name
	btn.action = me.Action

	opts := utils.OptionalArg(extraOpts)
	opts.WindowName = me.Text
	opts.ControlId = me.Id
	opts.Left = me.Pos.X
	opts.Top = me.Pos.Y
	if me.Disabled {
		opts.StyleInclude |= win32.WS_DISABLED
	}
	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := btn.Create(*opts)
	if me.Icon != 0 {
		btn.SetIcon(me.Icon)
	}
	assertNoErr(err)
	configControlSize(btn, me.Size)
	return btn
}

// ButtonObj implements Button.ButtonObj
func (this *ButtonObject) ButtonObj() *ButtonObject {
	return this
}

// SetDefault implements Button.SetDefault
func (this *ButtonObject) SetDefault() {
	style := this.GetStyle()
	style &^= WINDOW_STYLE(win32.BS_PUSHBUTTON)
	style |= WINDOW_STYLE(win32.BS_DEFPUSHBUTTON)
	SendMessage(this.Handle, win32.BM_SETSTYLE, WPARAM(style), 1)
	SendMessage(this.GetRootHandle(), win32.DM_SETDEFID, this.GetControlId(), 0)
}

// SetText implements Button.SetText
func (this *ButtonObject) SetText(text string) {
	win32.SetWindowText(this.Handle, win32.StrToPwstr(text))
}

// GetText implements Button.GetText
func (this *ButtonObject) GetText() string {
	text, _ := GetWindowText(this.Handle)
	return text
}

// SetIcon implements Button.SetIcon
func (this *ButtonObject) SetIcon(icon win32.HICON) {
	SendMessage(this.Handle, win32.BM_SETIMAGE,
		WPARAM(win32.IMAGE_ICON), icon)
}

// SetCommand implements Button.SetCommand
func (this *ButtonObject) SetCommand(command *Command) {
	this.command = command
	this.command.OnChange.AddListener(this.onCommandChange)
	this.updateFromCommand()
}

// GetCommand implements Button.GetCommand
func (this *ButtonObject) GetCommand() *Command {
	return this.command
}

// updateFromCommand updates the button based on the associated command
func (this *ButtonObject) updateFromCommand() {
	this.SetEnabled(!this.command.Disabled)
	if this.command.Text != "" {
		this.SetText(this.command.Text)
	}
}

// onCommandChange is the event handler for command changes
func (this *ButtonObject) onCommandChange(info *SimpleEventInfo) {
	this.updateFromCommand()
}

// GetWindowClass implements WindowSpi.GetWindowClass
func (this *ButtonObject) GetWindowClass() string {
	return "Button"
}

// Init implements Window.Init
func (this *ButtonObject) Init() {
	this.super.Init()
}

// OnReflectCommand implements WindowSpi.OnReflectCommand
func (this *ButtonObject) OnReflectCommand(msg *CommandMessage) {
	if msg.GetNotifyCode() == uint16(win32.BN_CLICKED) {
		this.RealObject.(ButtonSpi).OnClick()
	}
}

// OnReflectMessage implements WindowSpi.OnReflectMessage
func (this *ButtonObject) OnReflectMessage(msg *Message) {
	if msg.UMsg == win32.WM_CTLCOLORBTN {
		backColor := this.GetBackColor()
		if backColor != colors.Null && backColor != colors.Transparent {
			oldHbr := this.GetData(Data_BackColorBrush)
			if oldHbr != nil {
				win32.DeleteObject(oldHbr.(win32.HBRUSH))
			}
			hbr := win32.CreateSolidBrush(backColor.Win32Color())
			this.SetData(Data_BackColorBrush, hbr)
			msg.Result = hbr
			msg.Handled = true
		}
		return
	}
	this.super.OnReflectMessage(msg)
}

// OnReflectNotify implements WindowSpi.OnReflectNotify
func (this *ButtonObject) OnReflectNotify(info *NotifyMessage) {
	pNmhdr := info.GetNMHDR()
	foreColor := this.GetForeColor()
	if pNmhdr.Code == win32.NM_CUSTOMDRAW && foreColor != colors.Null {
		pNmcd := (*win32.NMCUSTOMDRAW)(unsafe.Pointer(pNmhdr))
		if pNmcd.DwDrawStage == win32.CDDS_PREPAINT {
			win32.SetTextColor(pNmcd.Hdc, foreColor.Win32Color())
			win32.SetBkMode(pNmcd.Hdc, win32.TRANSPARENT)
			format := win32.DT_CENTER | win32.DT_VCENTER | win32.DT_SINGLELINE
			pwszText := win32.StrToPwstr(this.GetText())
			win32.DrawText(pNmcd.Hdc, pwszText, -1, &pNmcd.Rc, format)
			//info.Result = win32.LRESULT(win32.CDRF_SKIPDEFAULT | win32.CDRF_NOTIFYPOSTPAINT)
			//info.Handled = true
			lResult := win32.LRESULT(win32.CDRF_SKIPDEFAULT | win32.CDRF_NOTIFYPOSTPAINT)
			info.SetHandledWithResult(lResult)
			return
		}
	}
	this.super.OnReflectNotify(info)
}

// OnClick implements ButtonSpi.OnClick
func (this *ButtonObject) OnClick() {
	if this.action != nil {
		this.action()
	}
	if this.command != nil {
		this.command.NotifyExecute()
	}
}

// GetPreferredSize implements Window.GetPreferredSize
func (this *ButtonObject) GetPreferredSize(int, int) (cx, cy int) {
	var sz win32.SIZE
	SendMessage(this.Handle, win32.BCM_GETIDEALSIZE,
		0, unsafe.Pointer(&sz))
	return int(sz.Cx), int(sz.Cy)
}

// Click implements Button.Click
func (this *ButtonObject) Click() {
	SendMessage(this.Handle, win32.BM_CLICK, 0, 0)
}
