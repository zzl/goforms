package forms

import (
	"github.com/zzl/goforms/drawing/colors"
	"log"
	"syscall"

	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing/gdi"
	"github.com/zzl/goforms/layouts"
)

// Container is an interface that represents a container window,
// which can hold other controls and manage their layout.
//
// todo: scrollable?
type Container interface {
	CustomWindow // the parent interface

	// Add adds one or more controls to the container.
	Add(control ...Control)

	// Remove removes the specified control from the container.
	Remove(ctrl Control)

	// SetLayout sets the layout for the container.
	SetLayout(layout Layout)

	// GetLayout returns the current layout of the container.
	GetLayout() Layout

	// UpdateLayout updates the layout of the container.
	UpdateLayout()

	// GetChildWindows returns direct child windows.
	GetChildWindows() []Window

	// GetDescendantWindows returns all descendant windows.
	GetDescendantWindows() []Window

	// GetChildControls returns  the child controls of the container.
	GetChildControls() []Control

	// GetControlByName returns the control with the specified name.
	GetControlByName(name string) Control

	// GetControlById returns the control with the specified ID.
	GetControlById(id int) Control

	// GetRootContainer returns the root container of the hierarchy.
	GetRootContainer() Container

	// UseDialogFont sets the container to use the dialog font.
	UseDialogFont()
}

// ContainerSpi is an interface that provides additional methods
// specific to implementing a Container.
type ContainerSpi interface {
	CustomWindowSpi
}

// ContainerInterface is a composition of ContainerWindow and ContainerSpi
type ContainerInterface interface {
	Container
	ContainerSpi
}

// ContainerObject implements the ContainerInterface
// It extends CustomWindowObject.
type ContainerObject struct {

	// CustomWindowObject is the parent struct.
	CustomWindowObject

	// super is the special pointer to the parent struct.
	super *CustomWindowObject

	Layout Layout // The layout of the container.

	childControls []Control // Slice containing child controls.

}

// ContextContainer is the default parent window
// for newly created controls if no explicit parent is specified.
var ContextContainer Container

// ContextContainerRestorer is the holder of the original ContextContainer
type ContextContainerRestorer struct {
	oriContainer Container
}

// Restore sets ContextContainer to the saved original one.
func (me ContextContainerRestorer) Restore() {
	ContextContainer = me.oriContainer
}

// SetContextContainer sets the specified container as current ContextContainer
// and returns ContextContainerRestorer that can be used to Restore the old one.
func SetContextContainer(container Container) ContextContainerRestorer {
	oriContainer := ContextContainer
	ContextContainer = container
	return ContextContainerRestorer{oriContainer: oriContainer}
}

// LayoutContainer is an layouts.Container implementation.
type LayoutContainer struct {
	containerObj *ContainerObject
}

// GetControlByName implements layouts.Container.GetControlByName.
func (me LayoutContainer) GetControlByName(name string) layouts.Control {
	return me.containerObj.GetControlByName(name)
}

// GetControls implements layouts.Container.GetControls.
func (me LayoutContainer) GetControls() []layouts.Control {
	ccs := me.containerObj.GetChildControls()
	count := len(ccs)
	lccs := make([]layouts.Control, count)
	for n := 0; n < count; n++ {
		lccs[n] = ccs[n]
	}
	return lccs
}

// GetClientSize implements layouts.Container.GetClientSize.
func (me LayoutContainer) GetClientSize() (cx, cy int) {
	return me.containerObj.GetClientSize()
}

// GetRootContainer implements Container.GetRootContainer.
func (this *ContainerObject) GetRootContainer() Container {
	return this.GetRootWindow().(Container)
}

var _containerClassRegstered = false

// EnsureClassRegistered implements WindowSpi.EnsureClassRegistered.
func (this *ContainerObject) EnsureClassRegistered() {
	if _containerClassRegstered {
		return
	}
	_, err := RegisterClass("goforms.container", nil, ClassOptions{
		BackgroundBrush: 0,
	})
	if err != nil {
		log.Fatal(err)
	}
	_containerClassRegstered = true
}

// EnsureCustomWndProc implements WindowSpi.EnsureCustomWndProc.
func (this *ContainerObject) EnsureCustomWndProc() {
	//nop
}

// GetWindowClass implements WindowSpi.GetWindowClass.
func (this *ContainerObject) GetWindowClass() string {
	return "goforms.container"
}

// PreCreate implements WindowSpi.PreCreate.
func (this *ContainerObject) PreCreate(opts *WindowOptions) {
	this.super.PreCreate(opts)
	hWndParent := opts.ParentHandle
	if this.GetFont() == nil && hWndParent != 0 {
		parentWin := GetWindow(hWndParent)
		font := parentWin.GetFont()
		if font != nil {
			font = font.CopyUnowned()
		}
		this.SetFont(font)
	}
}

// PreDispose implements WindowSpi.PreDispose.
func (this *ContainerObject) PreDispose() {
	childWins := this.GetChildWindows()
	for _, childWin := range childWins {
		childWin.AsWindowObject().callPreDispose()
	}
	this.SetData("Container.ChildWins", childWins)
	this.super.PreDispose()
}

// Init implements Window.Init
func (this *ContainerObject) Init() {
	this.super.Init()
}

// Dispose implements Window.Dispose.
func (this *ContainerObject) Dispose() {
	childWins := this.GetData("Container.ChildWins").([]Window)
	for _, childWin := range childWins {
		childWin.AsWindowObject().callDispose()
	}
	this.super.Dispose()
}

// UseDialogFont implements Container.UseDialogFont.
func (this *ContainerObject) UseDialogFont() {
	font := gdi.NewFontPt("MS Shell Dlg 2", 8)
	this.SetFont(font)
}

// GetChildControls implements Container.GetChildControls
func (this *ContainerObject) GetChildControls() []Control {
	return this.childControls
}

// GetControlByName implements Container.GetControlByName.
func (this *ContainerObject) GetControlByName(name string) Control {
	for _, c := range this.childControls {
		tName := c.GetName()
		if tName == "" {
			continue
		}
		if tName == name {
			return c
		}
	}
	return nil
}

// GetControlById implements Container.GetControlById.
func (this *ContainerObject) GetControlById(id int) Control {
	hWndItem, _ := win32.GetDlgItem(this.Handle, int32(id))
	if hWndItem != 0 {
		return GetWindow(hWndItem).(Control)
	}
	return nil
}

var _enumedDescendantWindows []Window
var _enumChildWindowsCallback = syscall.NewCallback(
	func(hWnd HWND, lparam LPARAM) win32.LRESULT {
		win, ok := windowMap[hWnd]
		if ok {
			_enumedDescendantWindows = append(_enumedDescendantWindows, win)
		}
		return 1
	})

// GetDescendantWindows implements Container.GetDescendantWindows.
func (this *ContainerObject) GetDescendantWindows() []Window {
	_enumedDescendantWindows = nil
	win32.EnumChildWindows(this.Handle, _enumChildWindowsCallback, 0)
	return _enumedDescendantWindows
}

// GetChildWindows implements Container.GetChildWindows.
func (this *ContainerObject) GetChildWindows() []Window {
	var childWins []Window
	childHwnds := GetChildHandles(this.Handle)
	for _, hWnd := range childHwnds {
		win, ok := windowMap[hWnd]
		if ok {
			childWins = append(childWins, win)
		}
	}
	return childWins
}

// GetLayout implements Container.GetLayout.
func (this *ContainerObject) GetLayout() Layout {
	return this.Layout
}

// SetLayout implements Container.SetLayout.
func (this *ContainerObject) SetLayout(layout Layout) {
	this.Layout = layout
	layout.SetContainer(LayoutContainer{this})
}

// UpdateLayout implements Container.UpdateLayout.
func (this *ContainerObject) UpdateLayout() {
	this.Layout.Update()
}

// SetBounds implements Window.SetBounds.
func (this *ContainerObject) SetBounds(left, top, width, height int) {
	var rc win32.RECT
	win32.GetWindowRect(this.Handle, &rc)
	if int(rc.Right-rc.Left) == width &&
		int(rc.Bottom-rc.Top) == height {
		return
	}
	this.super.SetBounds(left, top, width, height)
	if this.Layout != nil {
		this.Layout.SetBounds(0, 0, width, height)
	}
}

// Add implements Control.Add.
func (this *ContainerObject) Add(controls ...Control) {
	hWndSet := make(map[HWND]bool)
	for _, c := range this.childControls {
		hWndSet[c.GetHandle()] = true
	}
	for _, c := range controls {
		if c.GetHandle() == 0 { //?
			log.Panic("????")
		}
		if hWndSet[c.GetHandle()] { //already exist
			continue
		}
		c = c.GetRealObject().(Control)
		this.childControls = append(this.childControls, c)
	}
}

// Remove implements Control.Remove.
func (this *ContainerObject) Remove(ctrl Control) {
	var ctrls2 []Control
	hWnd := ctrl.GetHandle()
	for _, c := range this.childControls {
		if c.GetHandle() != hWnd {
			ctrls2 = append(ctrls2, c)
		}
	}
	this.childControls = ctrls2
	win32.DestroyWindow(hWnd)
	ctrl.Dispose()
}

// OnCtlColorStatic implements CustomWindowSpi.OnCtlColorStatic.
func (this *ContainerObject) OnCtlColorStatic(msg *Message) {
	backColor := this.GetBackColor()
	if backColor != colors.Null && !backColor.IsTransparent() {
		oldHbr := this.GetData(Data_BackColorBrush)
		if oldHbr != nil {
			win32.DeleteObject(oldHbr.(win32.HBRUSH))
		}
		hdc := msg.WParam
		win32.SetBkMode(hdc, win32.TRANSPARENT)
		hbr := win32.CreateSolidBrush(backColor.Win32Color())
		this.SetData(Data_BackColorBrush, hbr)
		msg.SetHandledWithResult(hbr)
	}
}
