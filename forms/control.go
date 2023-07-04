package forms

import (
	"github.com/zzl/goforms/framework/utils"
	"log"

	"github.com/zzl/go-win32api/v2/win32"
)

// Control is an interface that represents a graphical control window.
// It extends the Window and NameAware interfaces.
type Control interface {
	Window    //Control is a Window
	NameAware //Control is name aware

	// GetContainer returns the parent container of the control.
	GetContainer() Container

	// GetRootContainer returns the root container of the control.
	GetRootContainer() Container

	// CreateIn creates the control within a parent window.
	CreateIn(parent Window, extraOpts ...*WindowOptions) Control

	// GetControlSpecStyle returns the specific window styles for the control.
	GetControlSpecStyle() (include, exclude WINDOW_STYLE)

	// GetControlId returns the Win32 control identifier.
	GetControlId() uint16
}

// ControlSpi is an interface that provides additional methods
// specific to implementing a Control.
type ControlSpi interface {
	WindowSpi
}

// ControlInterface is a composition of Control and ControlSpi
type ControlInterface interface {
	Control
	ControlSpi
}

// ControlObject implements the ControlInterface
// It extends WindowObject.
type ControlObject struct {

	// WindowObject is the parent struct.
	WindowObject

	// super is a pointer to the parent struct.
	// It's used to ease the burden of remembering the
	// concrete parent struct and typing.
	//
	// The field name 'super' is special because virtual.Realize
	// will automatically set it up.
	super *WindowObject

	// NameAwareSupport is a NameAware implementation
	NameAwareSupport
}

// DefaultControlStyle is the default window style for controls.
var DefaultControlStyle WINDOW_STYLE = win32.WS_CHILD | win32.WS_VISIBLE | win32.WS_TABSTOP

// creatingControlMap maps Control windows with their control id.
//
// This map is used to temporarily associate Control windows with their
// control ID during the creation process before they are added to the windowMap.
var creatingControlMap map[int]WindowInterface = make(map[int]WindowInterface)

// autoControlIdBase is the lower bound of auto generated control ids.
const autoControlIdBase = 0x4000

// autoControlIdGen is a unique id generator used to auto generate control ids.
var autoControlIdGen = NewUidGen(autoControlIdBase, 1)

// Init implements Window.Init.
func (this *ControlObject) Init() {
	this.super.Init()
}

// PreDispose implements Window.PreDispose.
func (this *ControlObject) PreDispose() {
	id := this.GetControlId()
	if id >= autoControlIdBase && id <= 0xDFFF {
		autoControlIdGen.Recycle(int(id))
	}
	delete(creatingControlMap, int(id))
	this.super.PreDispose()
}

// GetRootContainer implements Control.GetRootContainer.
func (this *ControlObject) GetRootContainer() Container {
	return this.GetRootWindow().(Container)
}

// GetControlSpecStyle implements Control.GetControlSpecStyle.
func (this *ControlObject) GetControlSpecStyle() (include, exclude WINDOW_STYLE) {
	return 0, 0
}

// GetDefaultStyle implements WindowSpi.GetDefaultStyle.
func (this *ControlObject) GetDefaultStyle() WINDOW_STYLE {
	style := DefaultControlStyle
	incStyle, excStyle := this.RealObject.(Control).GetControlSpecStyle()
	style |= incStyle
	style &^= excStyle
	return style
}

// PostCreate implements WindowSpi.PostCreate.
func (this *ControlObject) PostCreate(opts *WindowOptions) {
	if container, ok := windowMap[this.GetParentHandle()].(Container); ok {
		container.Add(this)
	}
}

// createControlIn is a shared helper function to create control windows.
func createControlIn(parentWin Window, controlWin Window,
	extraOpts ...*WindowOptions) Control {
	opts := utils.OptionalArg(extraOpts)
	opts.ParentHandle = resolveParentHandle(parentWin)
	err := controlWin.Create(*opts)
	if err != nil {
		log.Fatal(err)
	}
	return controlWin.(Control)
}

// CreateIn implements Control.CreateIn.
func (this *ControlObject) CreateIn(parent Window, extraOpts ...*WindowOptions) Control {
	return createControlIn(parent, this.RealObject, extraOpts...)
}

// Create implements Window.Create.
func (this *ControlObject) Create(options WindowOptions) error {
	if options.ParentHandle == 0 {
		log.Println("Warning: Parent handle unspecified")
	}
	if options.ControlId == 0 {
		options.ControlId = uint16(autoControlIdGen.Gen())
	}
	creatingControlMap[int(options.ControlId)] = this.RealObject
	return this.super.Create(options)
}

// GetControlId implements Control.GetControlId.
func (this *ControlObject) GetControlId() uint16 {
	ret, errno := win32.GetDlgCtrlID(this.Handle)
	if ret == 0 {
		log.Println(errno.Error())
	}
	return uint16(ret)
}

// GetContainer implements Control.GetContainer.
func (this *ControlObject) GetContainer() Container {
	hWndParent, _ := win32.GetParent(this.Handle)
	parentWin := GetWindow(hWndParent)
	if parentWin == nil {
		return nil //?
	}
	return parentWin.(Container)
}
