package forms

import (
	"github.com/zzl/go-com/com"
	"github.com/zzl/goforms/framework/consts"
	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/framework/virtual"
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
)

type ToolBar interface {
	Control

	SetImageList(imageList *ImageList, imageListLarge *ImageList)
	SetLarge(large bool)
	IsLarge() bool

	AddItems(items []*ToolBarItem)
	BuildOverflowMenu() *PopupMenu

	GetItems() []*ToolBarItem
	GetItemIds() []int
	SetItemIds(ids []int)
	SetDefaultItemIds(ids []int)
	LoadSystemImages(idbId uint32) int

	EnableItem(id int, enabled bool)
}

type ToolBarSpi interface {
	ControlSpi
}

type ToolBarInterface interface {
	ToolBar
	ToolBarSpi
}

type ToolBarItem struct {
	Id            uint16
	Name          string
	Text          string
	Image         int
	Tooltip       string
	NoTextOnRight bool
	Action        Action
	Menu          *PopupMenu

	Disabled bool
	Command  *Command
}

type TbLabelStyle byte

const ID_SEPARATOR = 65534

const (
	TbLabelDefault   TbLabelStyle = 0 //default below
	TbLabelBelow     TbLabelStyle = 1
	TbLabelOnRight   TbLabelStyle = 2
	TbLabelInvisible TbLabelStyle = 3
)

type ToolBarObject struct {
	ControlObject
	super *ControlObject

	Flat          bool
	ButtonSize    Size
	LabelStyle    TbLabelStyle
	OnContextMenu SimpleEvent
	OnDblClick    SimpleEvent
	AllItems      []*ToolBarItem

	imageList      *ImageList
	imageListLarge *ImageList
	large          bool
	//iconOnlyIds    intsets.Sparse

	items     []*ToolBarItem
	idItemMap map[uint16]*ToolBarItem
	idGen     UidGen

	defaultItemIds []int
}

type NewToolBar struct {
	Parent     Container
	Name       string
	LabelStyle TbLabelStyle
}

func (me NewToolBar) Create(extraOpts ...*WindowOptions) ToolBar {
	toolBar := NewToolBarObject()
	toolBar.name = me.Name
	toolBar.LabelStyle = me.LabelStyle

	opts := utils.OptionalArg(extraOpts)
	opts.ParentHandle = resolveParentHandle(me.Parent)
	err := toolBar.Create(*opts)
	assertNoErr(err)

	return toolBar
}

func NewToolBarObject() *ToolBarObject {
	return virtual.New[ToolBarObject]()
}

func (this *ToolBarObject) LoadSystemImages(idbId uint32) int {
	ret, errno := win32.SendMessage(this.Handle,
		win32.TB_LOADIMAGES, win32.WPARAM(idbId), NegativeOne)
	_ = errno //?
	return int(ret)
}

func (this *ToolBarObject) PreProcessMsg(msg *win32.MSG) bool {
	if msg.Message == win32.WM_LBUTTONDBLCLK {
		x, y, _ := ParseMouseMsgParams(msg.WParam, msg.LParam)
		pt := win32.POINT{X: x, Y: y}
		ret, _ := SendMessage(this.Handle, win32.TB_HITTEST,
			0, unsafe.Pointer(&pt))
		win32.DispatchMessage(msg)
		if int32(ret) < 0 {
			this.OnDblClick.Fire(this, &SimpleEventInfo{})
		}
		return true
	}
	return false
}

func (this *ToolBarObject) Init() {
	this.super.Init()

	this.Flat = true
	this.idGen = UidGen{nextId: ToolbarGenIdStart, step: -1}
	this.idItemMap = make(map[uint16]*ToolBarItem)
}

func (this *ToolBarObject) GetWindowClass() string {
	return "ToolbarWindow32"
}

func (this *ToolBarObject) GetControlSpecStyle() (WINDOW_STYLE, WINDOW_STYLE) {
	var style WINDOW_STYLE
	if this.Flat {
		style |= WINDOW_STYLE(win32.TBSTYLE_FLAT)
	}
	style |= WINDOW_STYLE(win32.TBSTYLE_TOOLTIPS)
	style |= WINDOW_STYLE(win32.CCS_NOPARENTALIGN)
	style |= WINDOW_STYLE(win32.CCS_NODIVIDER)

	style |= WINDOW_STYLE(win32.TBSTYLE_TRANSPARENT) //?

	style |= win32.WS_CLIPSIBLINGS | win32.WS_CLIPCHILDREN //?

	//TBSTYLE_WRAPABLE?
	return style, win32.WS_TABSTOP
}

func toolBarSubclassProc(hWnd HWND, uMsg uint32, wParam WPARAM, lParam LPARAM,
	uIdSubclass uintptr, dwRefData uintptr) win32.LRESULT {
	if uMsg == win32.WM_MEASUREITEM {
		return GetWindow(hWnd).(*ToolBarObject).onMeasureItem(wParam, lParam)
	} else if uMsg == win32.WM_DRAWITEM {
		return GetWindow(hWnd).(*ToolBarObject).onDrawItem(wParam, lParam)
	}
	return win32.DefSubclassProc(hWnd, uMsg, wParam, lParam)
}

//todo:check if from menu
func (this *ToolBarObject) onMeasureItem(wParam WPARAM, lParam LPARAM) win32.LRESULT {
	lpmis := (*win32.MEASUREITEMSTRUCT)(unsafe.Pointer(lParam))
	id := lpmis.ItemID
	item := this.idItemMap[uint16(id)]
	hdc := win32.GetDC(this.Handle)

	text := item.Text
	if text == "" && item.Command != nil {
		text = item.Command.Text
	}
	wsz, _ := syscall.UTF16FromString(text)
	var size win32.SIZE
	win32.GetTextExtentPoint32(hdc, &wsz[0], int32(len(wsz)), &size)

	lpmis.ItemWidth = uint32(int(size.Cx) + this.imageList.Cx + 4 + 12)
	lpmis.ItemHeight = uint32(this.imageList.Cy + 3 + 3)

	win32.ReleaseDC(this.Handle, hdc)
	return 1
}

//todo:check if from menu
func (this *ToolBarObject) onDrawItem(wParam WPARAM, lParam LPARAM) win32.LRESULT {
	lpdis := (*win32.DRAWITEMSTRUCT)(unsafe.Pointer(lParam))
	id := lpdis.ItemID
	item := this.idItemMap[uint16(id)]
	text := item.Text
	if text == "" && item.Command != nil {
		text = item.Command.Text
	}
	wsz, _ := syscall.UTF16FromString(text)
	rc := lpdis.RcItem

	iml := this.imageList

	dtFlags := win32.DT_LEFT | win32.DT_VCENTER | win32.DT_SINGLELINE

	wszThemeClass, _ := syscall.UTF16FromString("Menu")
	hTheme := win32.OpenThemeData(this.Handle, &wszThemeClass[0])
	//hTheme := win32.HTHEME(0)

	var rcText = rc
	rcText.Left += int32(iml.Cx + 12)
	rcText.Top += 3

	selected := lpdis.ItemState&win32.ODS_SELECTED != 0
	if hTheme != 0 {
		var state win32.BARITEMSTATES
		if selected {
			state = win32.MBI_HOT
		} else {
			state = win32.MBI_NORMAL
		}

		clrBg := win32.GetSysColor(win32.COLOR_MENU)
		FillSolidRect(lpdis.HDC, &rc, clrBg)

		win32.DrawThemeBackground(hTheme, lpdis.HDC, int32(win32.MENU_BARITEM), int32(state), &rc, nil)
		win32.DrawThemeTextEx(hTheme, lpdis.HDC, int32(win32.MENU_BARITEM), int32(state),
			&wsz[0], int32(len(wsz)), dtFlags, &rcText, nil)

		win32.CloseThemeData(hTheme)
	} else {
		var textColor, bgColor win32.COLORREF
		if selected {
			textColor = win32.GetSysColor(win32.COLOR_HIGHLIGHTTEXT)
			bgColor = win32.GetSysColor(win32.COLOR_HIGHLIGHT)
		} else {
			textColor = win32.GetSysColor(win32.COLOR_MENUTEXT)
			bgColor = win32.GetSysColor(win32.COLOR_MENU)
		}
		FillSolidRect(lpdis.HDC, &rc, bgColor)
		oriTextColor := win32.SetTextColor(lpdis.HDC, textColor)
		win32.SetBkMode(lpdis.HDC, win32.TRANSPARENT)
		win32.DrawText(lpdis.HDC, &wsz[0], int32(len(wsz)), &rcText, dtFlags)
		win32.SetTextColor(lpdis.HDC, oriTextColor)
	}

	image := item.Image
	if image == 0 && item.Command != nil {
		image = item.Command.Image
	}

	win32.ImageList_DrawEx(iml.handle, int32(image), lpdis.HDC,
		rc.Left+4, rc.Top+3, int32(iml.Cx), int32(iml.Cy),
		win32.CLR_NONE_U, win32.CLR_DEFAULT_U, win32.ILD_NORMAL)

	return 1
}

func (this *ToolBarObject) Create(options WindowOptions) error {
	err := this.super.Create(options)
	if err != nil {
		log.Fatal(err)
	}

	//
	pSubclassProc := syscall.NewCallback(toolBarSubclassProc)
	bOk := win32.SetWindowSubclass(this.Handle, pSubclassProc, 0, 0)
	if bOk == win32.FALSE {
		println("?")
	}

	//
	SendMessage(this.Handle, win32.TB_SETEXTENDEDSTYLE,
		0, win32.TBSTYLE_EX_DRAWDDARROWS)

	SendMessage(this.Handle, win32.TB_BUTTONSTRUCTSIZE,
		unsafe.Sizeof(win32.TBBUTTON{}), 0)

	if this.imageList != nil {
		hIml := this.imageList.handle
		if this.large {
			hIml = this.imageListLarge.handle
		}
		SendMessage(this.Handle, win32.TB_SETIMAGELIST, 0, hIml)
	}

	if this.items != nil {
		this.fillItems(this.items)
	}

	labelStyle := this.LabelStyle
	this.LabelStyle = 4 //
	this.SetLabelStyle(labelStyle)

	//this.AddMsgFilter(this)

	return err
}

func (this *ToolBarObject) Dispose() {
	if this.imageList != nil {
		this.imageList.Dispose()
	}
	if this.imageListLarge != nil {
		this.imageListLarge.Dispose()
	}
	this.super.Dispose()
}

func (this *ToolBarObject) SetLabelStyle(labelStyle TbLabelStyle) {
	if this.LabelStyle == labelStyle {
		return
	}
	if labelStyle == TbLabelInvisible {
		SendMessage(this.Handle, win32.TB_SETMAXTEXTROWS, 0, 0)
	} else {
		SendMessage(this.Handle, win32.TB_SETMAXTEXTROWS, 1, 0)

		style := this.GetStyle()
		newStyle := style
		ret, _ := SendMessage(this.Handle, win32.TB_GETEXTENDEDSTYLE, 0, 0)
		extStyle := uint32(ret)
		newExtStyle := extStyle
		if labelStyle == TbLabelOnRight {
			if style&WINDOW_STYLE(win32.TBSTYLE_LIST) == 0 {
				newStyle |= WINDOW_STYLE(win32.TBSTYLE_LIST)
			}
			if extStyle&win32.TBSTYLE_EX_MIXEDBUTTONS == 0 {
				newExtStyle |= win32.TBSTYLE_EX_MIXEDBUTTONS
			}
		} else {
			if style&WINDOW_STYLE(win32.TBSTYLE_LIST) != 0 {
				newStyle &^= WINDOW_STYLE(win32.TBSTYLE_LIST)
			}
			if extStyle&win32.TBSTYLE_EX_MIXEDBUTTONS != 0 {
				newExtStyle &^= win32.TBSTYLE_EX_MIXEDBUTTONS
			}
		}
		if style != newStyle {
			win32.SetWindowLong(this.Handle, win32.GWL_STYLE, int32(newStyle))
		}
		if extStyle != newExtStyle {
			SendMessage(this.Handle, win32.TB_SETEXTENDEDSTYLE,
				0, newExtStyle)
		}
	}
	this.LabelStyle = labelStyle
	//
	SendMessage(this.Handle, win32.TB_AUTOSIZE, 0, 0)
	if this.IsVisible() {
		//todo: notify parent window about the resize..
		//GetWindow(this.GetParentHandle()).(WindowSpi).OnChildResized(this.RealObject)
	}
}

func (this *ToolBarObject) SetImageList(imageList *ImageList, imageListLarge *ImageList) {
	if this.imageList != nil {
		this.imageList.Dispose()
	}
	this.imageList = imageList
	if this.imageListLarge != nil {
		this.imageListLarge.Dispose()
	}
	this.imageListLarge = imageListLarge

	if this.Handle != 0 {
		SendMessage(this.Handle, win32.TB_SETIMAGELIST,
			0, imageList.handle)
	}
}

func (this *ToolBarObject) SetButtonSize(cx, cy int) {
	_, errno := SendMessage(this.Handle, win32.TB_SETBUTTONSIZE, 0,
		win32.MAKELONG(uint16(cx), uint16(cy)))
	if errno != win32.NO_ERROR {
		println("?")
	}
}

func (this *ToolBarObject) SetButtonWidth(minWdith, maxWidth int) {
	_, errno := SendMessage(this.Handle, win32.TB_SETBUTTONWIDTH, 0,
		win32.MAKELONG(uint16(minWdith), uint16(maxWidth)))
	if errno != win32.NO_ERROR {
		println("?")
	}
}

func (this *ToolBarObject) SetBitmapSize(cx, cy int) {
	_, errno := SendMessage(this.Handle, win32.TB_SETBITMAPSIZE, 0,
		win32.MAKELONG(uint16(cx), uint16(cy)))
	if errno != win32.NO_ERROR {
		println("?")
	}
}

func (this *ToolBarObject) SetPadding(cx, cy int) {
	_, errno := SendMessage(this.Handle, win32.TB_SETPADDING, 0,
		win32.MAKELONG(uint16(cx), uint16(cy)))
	if errno != win32.NO_ERROR {
		println("?")
	}
}

func (this *ToolBarObject) IsLarge() bool {
	return this.large
}

func (this *ToolBarObject) SetLarge(large bool) {
	this.large = large
	hIml := this.imageList.handle
	if large {
		hIml = this.imageListLarge.handle
	}
	if this.Handle != 0 {
		SendMessage(this.Handle, win32.TB_SETIMAGELIST, 0, hIml)
		//force correct size?
		if this.LabelStyle == TbLabelOnRight {
			this.SetLabelStyle(TbLabelBelow)
			this.SetLabelStyle(TbLabelOnRight)
		}
		//
		SendMessage(this.Handle, win32.TB_AUTOSIZE, 0, 0)
		if this.IsVisible() {
			//todo: notify parent window about the resize..
			//GetWindow(this.GetParentHandle()).(WindowSpi).OnChildResized(this.RealObject)
		}
	}
}

func (this *ToolBarObject) fillItems(items []*ToolBarItem) {
	var tbbs []win32.TBBUTTON
	bstrs := com.NewBStrs()
	defer bstrs.Dispose()

	for _, item := range items {
		var tbb win32.TBBUTTON
		text := item.Text
		if text == "-" {
			tbb.FsStyle = uint8(win32.BTNS_SEP)
			tbb.IBitmap = 7
		} else {
			image := item.Image

			command := item.Command
			if command != nil {
				text = command.GetNoPrefixText()
				if text == "" {
					text = item.Text
				}
				image = command.Image
				if image == 0 {
					image = item.Image
				}
			}
			if image == consts.Zero {
				image = 0
			}

			tbb.IBitmap = int32(image)
			tbb.IdCommand = int32(item.Id)
			tbb.FsState = uint8(win32.TBSTATE_ENABLED)

			if item.Disabled {
				tbb.FsState &^= uint8(win32.TBSTATE_ENABLED)
			}

			if command != nil {
				if command.Disabled {
					tbb.FsState &^= uint8(win32.TBSTATE_ENABLED)
				}
				if command.Checked {
					tbb.FsState |= uint8(win32.TBSTATE_CHECKED)
				}
			}

			tbb.FsStyle = uint8(win32.BTNS_BUTTON | win32.BTNS_NOPREFIX)

			if !item.NoTextOnRight {
				tbb.FsStyle |= uint8(win32.BTNS_SHOWTEXT)
			} else {
				//println("?")
			}

			if item.Menu != nil {
				if item.Action == nil {
					tbb.FsStyle |= uint8(win32.BTNS_WHOLEDROPDOWN)
				} else {
					tbb.FsStyle |= uint8(win32.BTNS_DROPDOWN)
				}
			}

			//pwsz, _ := syscall.UTF16PtrFromString(text)
			bstr := bstrs.Add(text)
			tbb.IString = bstr.Addr()
		}
		tbbs = append(tbbs, tbb)
	}
	SendMessage(this.Handle, win32.TB_ADDBUTTONS,
		len(tbbs), unsafe.Pointer(&tbbs[0]))

	SendMessage(this.Handle, win32.TB_AUTOSIZE, 0, 0)
}

func (this *ToolBarObject) Clear() {
	ret, _ := SendMessage(this.Handle, win32.TB_BUTTONCOUNT, 0, 0)
	count := int(ret)
	for n := count - 1; n >= 0; n-- {
		SendMessage(this.Handle, win32.TB_DELETEBUTTON, n, 0)
	}
	this.items = nil
}

func (this *ToolBarObject) SetItems(items []*ToolBarItem) {
	this.Clear()
	this.AddItems(items)
	SendMessage(this.Handle, win32.TB_AUTOSIZE, 0, 0)
	if this.IsVisible() {
		//todo: notify parent window about the resize..
		//GetWindow(this.GetParentHandle()).OnChildResized(this.RealObject)
	}
}

func (this *ToolBarObject) AddItems(items []*ToolBarItem) {
	//var pItems []*ToolBarItem
	//for n := range items {
	//	pItems = append(pItems, &items[n])
	//}
	this.processItems(items)
	this.items = append(this.items, items...)
	if this.Handle != 0 {
		this.fillItems(items)
	}
}

func (this *ToolBarObject) processItems(items []*ToolBarItem) {
	for _, item := range items {
		command := item.Command
		if item.Id == 0 {
			if command != nil && command.Id != 0 {
				item.Id = uint16(command.Id)
			} else {
				item.Id = uint16(this.idGen.Gen())
			}
		}
		this.idItemMap[item.Id] = item
		if command != nil {
			tItem := item
			command.OnChange.AddListener(func(eventInfo *SimpleEventInfo) {
				this.updateCommandItem(tItem)
			})
		}
	}
}

func (this *ToolBarObject) UpdateCommandItems() {
	for _, item := range this.items {
		if item.Command != nil {
			this.updateCommandItem(item)
		}
	}
}

func (this *ToolBarObject) updateCommandItem(item *ToolBarItem) {
	command := item.Command
	if command == nil {
		println("??")
		return
	}
	var tbi win32.TBBUTTONINFO
	tbi.CbSize = uint32(unsafe.Sizeof(tbi))
	tbi.DwMask = win32.TBIF_COMMAND | win32.TBIF_TEXT | win32.TBIF_IMAGE | win32.TBIF_STATE
	tbi.IdCommand = int32(item.Id)
	image := command.Image
	if image != 0 {
		if image == consts.Zero {
			image = 0
		}
		tbi.IImage = int32(image)
	} else {
		image = item.Image
	}
	tbi.IImage = int32(image)

	text := command.GetNoPrefixText()
	if text == "" {
		text = item.Text
	}
	pwsz, _ := syscall.UTF16PtrFromString(text)
	tbi.PszText = pwsz

	if !command.Disabled {
		tbi.FsState |= uint8(win32.TBSTATE_ENABLED)
	}
	if command.Checked {
		tbi.FsState |= uint8(win32.TBSTATE_CHECKED)
	}
	SendMessage(this.Handle, win32.TB_SETBUTTONINFO,
		item.Id, unsafe.Pointer(&tbi))
}

func (this *ToolBarObject) OnReflectNotify(msg *NotifyMessage) {
	this.super.OnReflectNotify(msg)
	pNmhdr := msg.GetNMHDR()
	code := pNmhdr.Code
	if code == win32.TBN_DROPDOWN {
		pNmtb := (*win32.NMTOOLBAR)(unsafe.Pointer(pNmhdr))
		this.onDropdown(pNmtb)
	} else if code == win32.TBN_GETINFOTIP {
		pNmgit := (*win32.NMTBGETINFOTIP)(unsafe.Pointer(pNmhdr))
		this.onGetInfoTip(pNmgit)
	} else if code == win32.NM_RCLICK {
		this.OnContextMenu.Fire(this, &SimpleEventInfo{})
	} else if code == win32.TBN_BEGINADJUST {
		println("X")
	} else if code == win32.TBN_QUERYINSERT {
		msg.Result = 1
		msg.Handled = true
	}
}

func (this *ToolBarObject) onDropdown(pNmtb *win32.NMTOOLBAR) {
	id := uint16(pNmtb.IItem)
	item := this.idItemMap[id]
	var rc win32.RECT
	SendMessage(this.Handle, win32.TB_GETRECT,
		id, unsafe.Pointer(&rc))
	win32.MapWindowPoints(this.Handle, 0, (*win32.POINT)(unsafe.Pointer(&rc)), 2)
	item.Menu.Show(int(rc.Left), int(rc.Bottom), this.Handle)
}

func (this *ToolBarObject) onGetInfoTip(pNmgit *win32.NMTBGETINFOTIP) {
	id := uint16(pNmgit.IItem)
	item := this.idItemMap[id]
	tipText := item.Tooltip
	if tipText == "" {
		tipText = item.Text
	}
	command := item.Command
	if command != nil {
		tipText = command.Tooltip
		if tipText == "" {
			tipText = command.GetNoPrefixText()
		}
	}
	wsz, _ := syscall.UTF16FromString(tipText)
	cBuf := (*[1 << 29]uint16)(unsafe.Pointer(pNmgit.PszText))
	copy(cBuf[:pNmgit.CchTextMax], wsz)
}

func (this *ToolBarObject) GetPreferredSize(cxMax int, cyMax int) (int, int) {
	var size win32.SIZE
	SendMessage(this.Handle, win32.TB_GETMAXSIZE,
		0, unsafe.Pointer(&size))
	return int(size.Cx), int(size.Cy)
}

func (this *ToolBarObject) OnReflectCommand(msg *CommandMessage) {
	id := msg.GetCmdId()
	item := this.idItemMap[id]
	if item.Command != nil {
		item.Command.NotifyExecute()
	}
	if item.Action != nil {
		item.Action()
	}
}

func (this *ToolBarObject) OnParentResized() {
	this.super.OnParentResized()
	//?
	SendMessage(this.Handle, win32.TB_AUTOSIZE, 0, 0)
}

func (this *ToolBarObject) BuildOverflowMenu() *PopupMenu {
	count := len(this.items)
	ppm := NewPopupMenu()

	var rcClient win32.RECT
	win32.GetClientRect(this.Handle, &rcClient)

	for n := 0; n < count; n++ {
		var tbb win32.TBBUTTON
		ret, errno := SendMessage(this.Handle, win32.TB_GETBUTTON,
			n, unsafe.Pointer(&tbb))
		if ret == 0 {
			log.Fatal(errno)
		}
		if tbb.FsState&uint8(win32.TBSTATE_HIDDEN) != 0 {
			continue
		}
		var rcButton win32.RECT
		ret, errno = SendMessage(this.Handle, win32.TB_GETITEMRECT,
			n, unsafe.Pointer(&rcButton))
		if ret == 0 {
			log.Fatal(errno)
		}
		enabled := tbb.FsState&uint8(win32.TBSTATE_ENABLED) != 0
		if rcButton.Right <= rcClient.Right && rcButton.Bottom <= rcClient.Bottom {
			continue
		}
		if tbb.FsStyle&uint8(win32.BTNS_SEP) != 0 {
			if ppm.GetItemCount() > 0 {
				ppm.AddSeparator()
			}
			continue
		}
		tbItem := this.items[n]
		menuItem := MenuItem{
			Id:       tbItem.Id,
			Text:     tbItem.Text,
			Image:    tbItem.Image,
			Disabled: !enabled,
			Command:  tbItem.Command,
		}
		if tbItem.Menu != nil {
			menuItem.subMenuHandle = tbItem.Menu.Handle
			for id, item := range tbItem.Menu.idItemMap {
				if _, ok := ppm.idItemMap[id]; !ok {
					ppm.idItemMap[id] = item
				} else {
					log.Println("id conflict..")
				}
			}
		}
		ppm.AddItem(menuItem)
	}
	return ppm
}

func (this *ToolBarObject) LoadImages(typeId uint32) int {
	ret, errno := SendMessage(this.Handle, win32.TB_LOADIMAGES,
		typeId, NegativeOne)
	_ = errno
	return int(ret)
}

func (this *ToolBarObject) Customize() {
	//todo..
}

func (this *ToolBarObject) GetItems() []*ToolBarItem {
	return this.items
}

func (this *ToolBarObject) GetItemIds() []int {
	var ids []int
	for _, item := range this.items {
		id := int(item.Id)
		if id == 0 && item.Command != nil {
			id = item.Command.Id
		}
		if id == ID_SEPARATOR {
			id = 0
		}
		ids = append(ids, id)
	}
	return ids
}

func (this *ToolBarObject) SetItemIds(ids []int) {
	items := this.buildItemsFromIds(ids)
	this.SetItems(items)
}

func (this *ToolBarObject) buildItemsFromIds(ids []int) []*ToolBarItem {
	idItemMap := make(map[int]*ToolBarItem)
	for _, item := range this.AllItems {
		id := int(item.Id)
		if id == 0 && item.Command != nil {
			id = item.Command.Id
		}
		if id != 0 {
			idItemMap[id] = item
		}
	}
	var items []*ToolBarItem
	for _, id := range ids {
		if id == 0 || id == ID_SEPARATOR {
			items = append(items, &ToolBarItem{Id: ID_SEPARATOR, Text: "-"})
			continue
		}
		item, ok := idItemMap[id]
		if ok {
			items = append(items, item)
		}
	}
	return items
}

func (this *ToolBarObject) SetDefaultItemIds(ids []int) {
	this.defaultItemIds = ids
}

func (this *ToolBarObject) EnableItem(id int, enabled bool) {
	var nEnabled win32.LPARAM
	if enabled {
		nEnabled = win32.LPARAM(win32.TRUE)
	}
	SendMessage(this.Handle, win32.TB_ENABLEBUTTON, id, nEnabled)
}
