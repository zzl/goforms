package forms

import (
	"github.com/zzl/goforms/framework/consts"
	"log"
	"syscall"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing"
)

type MenuItem struct {
	Id       uint16
	Name     string
	Text     string
	Desc     string //description
	Image    int    //default 0 means null. use zero to specify 0.
	Radio    bool
	Disabled bool
	Checked  bool
	Action   Action
	SubItems []*MenuItem

	Command *Command

	ownerMenuHandle win32.HMENU
	subMenuHandle   win32.HMENU
	hBitmap         win32.HBITMAP
}

type Menu struct {
	Handle     win32.HMENU
	OnSelected func(item *MenuItem)

	ImageList *ImageList
	OwnerDraw bool

	items []*MenuItem
	idGen UidGen

	idItemMap   map[uint16]*MenuItem
	nameItemMap map[string]*MenuItem

	notifyByPos bool
}

type PopupMenu struct {
	Menu
}

// action not working?
type MenuBar struct {
	Menu
	boundWindowHandle HWND
}

var menuMap map[win32.HMENU]*Menu

func putMenuInMap(menu *Menu) {
	if menuMap == nil {
		menuMap = make(map[win32.HMENU]*Menu)
	}
	menuMap[menu.Handle] = menu
}

func GetMenuObject(hMenu win32.HMENU) *Menu {
	return menuMap[hMenu]
}

func (this *Menu) Init() {
	//
}

//create/destroy pair?

func (this *Menu) Dispose() {
	if this.Handle == 0 {
		return
	}
	this.disposeItems(this.items)
	win32.DestroyMenu(this.Handle)
	delete(menuMap, this.Handle)
	this.Handle = 0
}

func (this *Menu) disposeItem(item *MenuItem) {
	if item.SubItems != nil {
		this.disposeItems(item.SubItems)
	}
	if item.hBitmap != 0 {
		win32.DeleteObject(win32.HGDIOBJ(item.hBitmap))
		item.hBitmap = 0
	}
}

func (this *Menu) disposeItems(items []*MenuItem) {
	for _, item := range items {
		this.disposeItem(item)
	}
}

func (this *Menu) populateItems(hMenu win32.HMENU, items []*MenuItem) {

	if this.notifyByPos {
		var mi win32.MENUINFO
		mi.CbSize = uint32(unsafe.Sizeof(mi))
		mi.FMask = win32.MIM_STYLE
		mi.DwStyle = win32.MNS_NOTIFYBYPOS
		bOk, errno := win32.SetMenuInfo(hMenu, &mi)
		if bOk == win32.FALSE {
			log.Fatal(errno)
		}
	}

	var mii win32.MENUITEMINFO
	mii.CbSize = uint32(unsafe.Sizeof(mii))

	oriCount, _ := win32.GetMenuItemCount(hMenu)

	for n, item := range items {
		item.ownerMenuHandle = hMenu
		mii.FMask = win32.MIIM_FTYPE
		text := item.Text
		command := item.Command
		if command != nil {
			text = command.Text
			if len(command.ShortcutKeys) > 0 {
				text += "\t" + command.ShortcutKeys[0].String()
			}
		}
		if text == "-" {
			mii.FType = win32.MFT_SEPARATOR
		} else {
			mii.FType = win32.MFT_STRING
			mii.FMask |= win32.MIIM_STRING | win32.MIIM_ID | win32.MIIM_DATA
			if item.Radio {
				mii.FType = win32.MFT_STRING | win32.MFT_RADIOCHECK
			}
			mii.WID = uint32(item.Id)
			mii.DwItemData = uintptr(this.Handle)
			mii.DwTypeData, _ = syscall.UTF16PtrFromString(text)
		}
		image := item.Image
		disabled := item.Disabled
		checked := item.Checked
		if command != nil {
			disabled = command.Disabled
			checked = command.Checked
			if command.Image != 0 {
				image = command.Image
			}
		}
		if disabled || checked {
			mii.FMask |= win32.MIIM_STATE
			if disabled {
				mii.FState |= win32.MFS_DISABLED
			}
			if checked {
				mii.FState |= win32.MFS_CHECKED
			}
		}
		if image != 0 && this.ImageList != nil {
			if image == consts.Zero {
				image = 0
			}
			mii.FMask |= win32.MIIM_BITMAP
			hIcon := this.ImageList.GetHIcon(image)
			item.hBitmap = drawing.HIconToHBitmap(hIcon)
			mii.HbmpItem = item.hBitmap
		}

		//
		if item.SubItems != nil {
			hSubMenu, _ := win32.CreateMenu()
			item.subMenuHandle = hSubMenu
			this.populateItems(hSubMenu, item.SubItems)
		}
		if item.subMenuHandle != 0 {
			mii.FMask |= win32.MIIM_SUBMENU
			mii.HSubMenu = item.subMenuHandle
		}
		if this.OwnerDraw {
			if mii.FType&win32.MFT_SEPARATOR == 0 {
				mii.FType |= win32.MFT_OWNERDRAW
			}
		}
		ok, errno := win32.InsertMenuItem(hMenu,
			uint32(oriCount)+uint32(n), win32.TRUE, &mii)
		if ok == win32.FALSE {
			log.Fatal(errno)
		}
	}
}

func (this *Menu) GetItem(id uint16) *MenuItem {
	return this.idItemMap[id]
}

func (this *MenuBar) GetItemByName(name string) *MenuItem {
	return this.nameItemMap[name]
}

func (this *Menu) GetSubMenuItem(hSubMenu win32.HMENU) *MenuItem {
	for _, item := range this.items {
		if item.subMenuHandle == hSubMenu {
			return item
		}
	}
	return nil
}

func (this *Menu) SetItemText(id int, text string) {
	var mii win32.MENUITEMINFO
	mii.CbSize = uint32(unsafe.Sizeof(mii))
	mii.FMask = win32.MIIM_STRING
	mii.FType = win32.MFT_STRING
	pwszText, _ := syscall.UTF16PtrFromString(text)
	mii.DwTypeData = pwszText
	bOk, errno := win32.SetMenuItemInfo(this.Handle, uint32(id), win32.FALSE, &mii)
	if bOk == win32.FALSE {
		log.Fatal(errno)
	}
}

func (this *Menu) AddSeparator() {
	this.AddItem(MenuItem{Text: "-"})
}

func (this *Menu) AddItem(item MenuItem) {
	this.AddItems([]MenuItem{item})
}

func (this *Menu) AddItems(items []MenuItem) {
	var pItems []*MenuItem
	for n := range items {
		pItems = append(pItems, &items[n])
	}
	this.processItems(pItems)
	this.items = append(this.items, pItems...)
	if this.Handle != 0 {
		this.populateItems(this.Handle, pItems)
	}
}

func (this *Menu) GetItemCount() int {
	return len(this.items)
}

func (this *Menu) processItems(items []*MenuItem) {
	for _, item := range items {
		command := item.Command
		if item.Text == "-" {
			//nop
		} else if item.SubItems != nil {
			this.processItems(item.SubItems)
		} else {
			if item.Id == 0 {
				if command != nil && command.Id != 0 {
					item.Id = uint16(command.Id)
				} else {
					item.Id = uint16(this.idGen.Gen())
				}
			}
			this.idItemMap[item.Id] = item
			if item.Name != "" {
				this.nameItemMap[item.Name] = item
			}
		}
		if command != nil {
			tItem := item
			command.OnChange.AddListener(func(info *SimpleEventInfo) {
				this.updateCommandItem(tItem)
			})
		}
	}
}

func (this *Menu) updateCommandItem(item *MenuItem) {
	command := item.Command
	if command == nil {
		println("??")
		return
	}
	hMenu := item.ownerMenuHandle
	if hMenu == 0 {
		return
	}
	var mii win32.MENUITEMINFO
	mii.CbSize = uint32(unsafe.Sizeof(mii))
	mii.FMask = win32.MIIM_STRING | win32.MIIM_STATE
	text := command.Text
	if len(command.ShortcutKeys) > 0 {
		text += "\t" + command.ShortcutKeys[0].String()
	}
	mii.DwTypeData, _ = syscall.UTF16PtrFromString(text)

	if item.hBitmap != 0 {
		win32.DeleteObject(win32.HGDIOBJ(item.hBitmap))
		item.hBitmap = 0
	}
	image := item.Image
	if command.Image != 0 {
		image = command.Image
	}
	if image != 0 {
		if image == consts.Zero {
			image = 0
		}
		mii.FMask |= win32.MIIM_BITMAP
		hIcon := this.ImageList.GetHIcon(image)
		item.hBitmap = drawing.HIconToHBitmap(hIcon)
		mii.HbmpItem = item.hBitmap
	}

	if command.Disabled {
		mii.FState |= win32.MFS_DISABLED
	}
	if command.Checked {
		mii.FState |= win32.MFS_CHECKED
	}
	win32.SetMenuItemInfo(hMenu, uint32(item.Id), win32.FALSE, &mii)
}

func NewPopupMenu() *PopupMenu {
	menu := &PopupMenu{
		Menu{
			idGen:       UidGen{nextId: -1, step: -1},
			idItemMap:   make(map[uint16]*MenuItem),
			nameItemMap: make(map[string]*MenuItem),
			notifyByPos: false,
		},
	}
	return menu
}

func (this *PopupMenu) Create() {
	hMenu, _ := win32.CreatePopupMenu()
	this.Handle = hMenu
	this.populateItems(hMenu, this.items)
	putMenuInMap(&this.Menu)
}

func (this *PopupMenu) Show(x, y int, hWndOwner HWND) (uint16, string) {
	if this.Handle == 0 {
		log.Println("Menu not created?")
	}
	retVal, _ := win32.TrackPopupMenu(this.Handle,
		win32.TPM_RETURNCMD,
		int32(x), int32(y), 0, hWndOwner, nil)
	if retVal == 0 { //canceled..
		return 0, ""
	}
	id := uint16(retVal)
	item := this.GetItem(id)
	if item == nil {
		//?log.Fatal("menu item not found.. #" + strconv.Itoa(int(id)))
	} else {
		if item.Action != nil {
			item.Action()
		}
		if item.Command != nil {
			item.Command.NotifyExecute()
		}
	}
	if !this.idGen.IsGenerated(int(id)) {
		SendMessage(hWndOwner, win32.WM_COMMAND, id, 0)
	}
	//?
	name := this.GetItem(id).Name
	return id, name
}

func NewMenuBar(create bool) *MenuBar {
	menu := &MenuBar{
		Menu: Menu{
			idGen:       UidGen{nextId: MenuGenIdStart, step: -1},
			idItemMap:   make(map[uint16]*MenuItem),
			nameItemMap: make(map[string]*MenuItem),
			notifyByPos: true,
		},
	}
	if create {
		menu.Create()
	}
	return menu
}

func (this *MenuBar) Create() {
	hMenu, _ := win32.CreateMenu()
	this.Handle = hMenu
	this.populateItems(hMenu, this.items)
	putMenuInMap(&this.Menu)
}

func (this *MenuBar) BindTo(topWin TopWindow) {
	win32.SetMenu(topWin.GetHandle(), this.Handle)
	this.boundWindowHandle = topWin.GetHandle()
}

func (this *MenuBar) onMenuCommand(wParam WPARAM, lParam LPARAM) bool {
	index := int(wParam)
	hMenu := win32.HMENU(lParam)
	var mii win32.MENUITEMINFO
	mii.CbSize = uint32(unsafe.Sizeof(mii))
	mii.FMask = win32.MIIM_ID | win32.MIIM_DATA
	win32.GetMenuItemInfo(hMenu, uint32(index), 1, &mii)
	if mii.DwItemData != uintptr(this.Handle) {
		println("not for menubar..")
		return false
	}
	id := uint16(mii.WID)
	item := this.idItemMap[id]
	if item.Action != nil {
		item.Action()
	}
	if item.Command != nil {
		item.Command.NotifyExecute()
	}
	if this.OnSelected != nil {
		this.OnSelected(item)
	}
	SendMessage(this.boundWindowHandle, win32.WM_COMMAND, id, 0)
	return true
}

func (this *MenuBar) CheckItem(id uint16, checked bool) {
	item := this.GetItem(id)
	var uCheck uint32
	if checked {
		uCheck = uint32(win32.MF_CHECKED)
	} else {
		uCheck = uint32(win32.MF_UNCHECKED)
	}
	ret := win32.CheckMenuItem(item.ownerMenuHandle, uint32(id), uCheck)
	if int32(ret) == -1 {
		println("??")
	}
}

func (this *MenuBar) CheckItemByName(name string, checked bool) {
	item := this.GetItemByName(name)
	this.CheckItem(item.Id, checked)
}
