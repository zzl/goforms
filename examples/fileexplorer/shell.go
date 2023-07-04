package main

import (
	"encoding/binary"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/forms"
	"log"
	"slices"
	"sort"
	"syscall"
	"time"
	"unsafe"
)

var cat = syscall.GUID{0xB725F130, 0x47EF, 0x101A,
	[8]byte{0xA5, 0xF1, 0x02, 0x60, 0x8C, 0x9E, 0xEB, 0xAC}}

var PKEY_Size = win32.PROPERTYKEY{Fmtid: cat, Pid: 12}
var PKEY_DateModified = win32.PROPERTYKEY{Fmtid: cat, Pid: 14}
var PKEY_ItemTypeText = win32.PROPERTYKEY{Fmtid: cat, Pid: 4}
var PKEY_ItemNameDisplay = win32.PROPERTYKEY{Fmtid: cat, Pid: 10}

type Shell struct {
	pIml       *win32.IImageList
	psfDesktop *win32.IShellFolder
}

func (this *Shell) Init() {
	var hr win32.HRESULT

	hr = win32.SHGetImageList(int32(win32.SHIL_SMALL),
		&win32.IID_IImageList, unsafe.Pointer(&this.pIml))
	win32.ASSERT_SUCCEEDED(hr)

	hr = win32.SHGetDesktopFolder(&this.psfDesktop)
	win32.ASSERT_SUCCEEDED(hr)
}

func (this *Shell) Dispose() {
	this.psfDesktop.Release()
	this.pIml.Release()
}

func (this *Shell) GetImageList() win32.HIMAGELIST {
	hIml := (win32.HIMAGELIST)(unsafe.Pointer(this.pIml))
	return hIml
}

func (this *Shell) GetDesktop() *ShellItem {
	var pidl *win32.ITEMIDLIST

	hr := win32.SHGetKnownFolderIDList(&win32.FOLDERID_Desktop,
		uint32(win32.KF_FLAG_DEFAULT), 0, &pidl)
	win32.ASSERT_SUCCEEDED(hr)

	//var psf *win32.IShellFolder
	//hr = win32.SHBindToObject(nil, pidl, nil,
	//	&win32.IID_IShellFolder, unsafe.Pointer(&psf))
	//win32.ASSERT_SUCCEEDED(hr)

	return NewShellItem(nil, this.psfDesktop, pidl)
}

func (this *Shell) BrowserForFolder() *ShellItem {
	var bi win32.BROWSEINFO
	bi.HwndOwner = forms.HWndActive
	szDisplayName := make([]uint16, win32.MAX_PATH)
	bi.PszDisplayName = &szDisplayName[0]
	bi.LpszTitle = win32.StrToPwstr("Choose folder")
	bi.UlFlags = win32.BIF_NEWDIALOGSTYLE
	pidl := win32.SHBrowseForFolder(&bi)
	if pidl == nil {
		return nil
	}
	defer win32.CoTaskMemFree(unsafe.Pointer(pidl))
	return this.NewShellItemFromAbsPidl(pidl)
}

func (this *Shell) NewShellItemFromAbsPidl(pidl *win32.ITEMIDLIST) *ShellItem {

	parentItem := this.GetDesktop()
	psfParent := this.psfDesktop

	pidlBts := unsafe.Slice((*byte)(unsafe.Pointer(pidl)), 64*1024)
	var item *ShellItem
	var psfs []*win32.IShellFolder
	for {
		cb := int(pidl.Mkid.Cb)
		if cb == 0 {
			break
		}
		id := pidlBts[2:cb]
		relPidl := makeRelPidl(id)
		//println(len(id))

		item = NewShellItem(parentItem, psfParent, relPidl)

		parentItem = item

		var psf *win32.IShellFolder
		hr := psfParent.BindToObject(relPidl, nil, &win32.IID_IShellFolder, unsafe.Pointer(&psf))
		if win32.SUCCEEDED(hr) {
			psfParent = psf
			psfs = append(psfs, psf)
		} else {
			psfParent = nil //?
		}

		pidlBts = pidlBts[cb:]
		pidl = (*win32.ITEMIDLIST)(unsafe.Pointer(&pidlBts[0]))
	}
	for _, psf := range psfs {
		psf.Release()
	}
	println("?")
	return item
}

type ShellItem struct {
	parent *ShellItem

	id []byte

	Name string
	Icon int

	ModTime time.Time
	Type    string
	Size    int

	IsFolder bool
	HasChild bool

	_relPidl *win32.ITEMIDLIST //transient
}

func makeRelPidl(id []byte) *win32.ITEMIDLIST {
	if id == nil {
		bts00 := []byte{0, 0}
		return (*win32.ITEMIDLIST)(unsafe.Pointer(&bts00[0]))
	} else {
		cb := uint16(len(id) + 2)
		bts := make([]byte, cb+2)
		*(*uint16)(unsafe.Pointer(&bts[0])) = cb
		copy(bts[2:], id)
		return (*win32.ITEMIDLIST)(unsafe.Pointer(&bts[0]))
	}
}

// relative pidl
func NewShellItem(parentItem *ShellItem, psfParent *win32.IShellFolder,
	pidl *win32.ITEMIDLIST) *ShellItem {

	cb := pidl.Mkid.Cb
	var id []byte
	if cb != 0 {
		id = make([]byte, cb-2)
		copy(id, unsafe.Slice((*byte)(&pidl.Mkid.AbID[0]), cb-2))
	} else {
		id = nil
	}
	item := &ShellItem{
		id:       id,
		parent:   parentItem,
		_relPidl: pidl,
	}
	item.loadInfo(psfParent)
	return item
}

// no need to free
func (this *ShellItem) getAbsPidl() *win32.ITEMIDLIST {
	if this.id == nil {
		var bts00 = []byte{0, 0}
		return (*win32.ITEMIDLIST)(unsafe.Pointer(&bts00[0]))
	}
	var bts []byte
	pathItems := this.GetPathItems()
	//for n := len(pathItems) - 2; n >= 0; n-- {
	for n := 1; n < len(pathItems); n++ {
		item := pathItems[n]
		bts = binary.LittleEndian.AppendUint16(bts, uint16(len(item.id)+2))
		bts = append(bts, item.id...)
	}
	bts = append(bts, 0)
	bts = append(bts, 0)
	return (*win32.ITEMIDLIST)(unsafe.Pointer(&bts[0]))
}

func (this *ShellItem) GetPathItems() []*ShellItem {
	var items []*ShellItem
	item := this
	for item != nil {
		items = append(items, item)
		item = item.parent
	}
	slices.Reverse(items)
	return items
}

func (this *ShellItem) getPidl() *win32.ITEMIDLIST {
	return makeRelPidl(this.id)
}

func (this *ShellItem) loadInfo(psfParent *win32.IShellFolder) {

	var hr win32.HRESULT

	pidl := this._relPidl

	//
	var attrs uint32
	attrs = uint32(win32.SFGAO_FOLDER | win32.SFGAO_STREAM)
	hr = psfParent.GetAttributesOf(1, &pidl, &attrs)
	win32.ASSERT_SUCCEEDED(hr)
	this.IsFolder = attrs&uint32(win32.SFGAO_FOLDER) != 0 &&
		attrs&uint32(win32.SFGAO_STREAM) == 0

	if this.IsFolder {
		attrs = uint32(win32.SFGAO_HASSUBFOLDER)
		hr = psfParent.GetAttributesOf(1, &pidl, &attrs)
		win32.ASSERT_SUCCEEDED(hr)
		this.HasChild = attrs&uint32(win32.SFGAO_HASSUBFOLDER) != 0
	}

	//
	absPidl := this.getAbsPidl()
	var sfi win32.SHFILEINFO
	flags := win32.SHGFI_PIDL | win32.SHGFI_SMALLICON | win32.SHGFI_SYSICONINDEX
	win32.SHGetFileInfo(win32.PWSTR(unsafe.Pointer(absPidl)), 0,
		&sfi, uint32(unsafe.Sizeof(sfi)), flags)
	this.Icon = int(sfi.IIcon)

	//
	var psi *win32.IShellItem2
	hr = win32.SHCreateItemWithParent(nil, psfParent, this._relPidl,
		&win32.IID_IShellItem2, unsafe.Pointer(&psi))
	win32.ASSERT_SUCCEEDED(hr)
	defer psi.Release()

	var pwsz win32.PWSTR

	psi.GetString(&PKEY_ItemNameDisplay, &pwsz)
	if pwsz != nil {
		this.Name = win32.PwstrToStr(pwsz)
		win32.CoTaskMemFree(unsafe.Pointer(pwsz))
	}

	var size uint64
	hr = psi.GetUInt64(&PKEY_Size, &size)
	if win32.SUCCEEDED(hr) {
		this.Size = int(size)
	} else {
		this.Size = -1
	}

	var ft win32.FILETIME
	hr = psi.GetFileTime(&PKEY_DateModified, &ft)
	if win32.SUCCEEDED(hr) {
		t := time.Unix(0, (*syscall.Filetime)(unsafe.Pointer(&ft)).Nanoseconds())
		this.ModTime = t
	}

	psi.GetString(&PKEY_ItemTypeText, &pwsz)
	if pwsz != nil {
		this.Type = win32.PwstrToStr(pwsz)
		win32.CoTaskMemFree(unsafe.Pointer(pwsz))
	}

}

func (this *ShellItem) GetChildren(includeNonFolder bool) []*ShellItem {
	pidl := this.getAbsPidl()
	var psf *win32.IShellFolder
	hr := win32.SHBindToObject(nil, pidl, nil,
		&win32.IID_IShellFolder, unsafe.Pointer(&psf))
	win32.ASSERT_SUCCEEDED(hr)

	var pEnumIdList *win32.IEnumIDList
	flags := uint32(win32.SHCONTF_FOLDERS)
	if includeNonFolder {
		flags |= uint32(win32.SHCONTF_NONFOLDERS)
	}
	//flags |= uint32(win32.SHCONTF_FASTITEMS) //?
	hr = psf.EnumObjects(forms.HWndActive, flags, &pEnumIdList)
	if hr == win32.S_FALSE {
		return nil
	}
	if win32.FAILED(hr) { //??
		log.Println(win32.HRESULT_ToString(hr))
		return nil
	}
	//win32.ASSERT_SUCCEEDED(hr)
	defer pEnumIdList.Release()

	var items []*ShellItem

	var celtFetched uint32
	for {
		var pidl *win32.ITEMIDLIST
		hr = pEnumIdList.Next(1, &pidl, &celtFetched)
		if !win32.SUCCEEDED(hr) || celtFetched != 1 {
			break
		}

		//
		if !includeNonFolder {
			var attrs uint32
			attrs = uint32(win32.SFGAO_STREAM)
			hr = psf.GetAttributesOf(1, &pidl, &attrs)
			win32.ASSERT_SUCCEEDED(hr)
			if (attrs & uint32(win32.SFGAO_STREAM)) != 0 {
				continue
			}
		}

		//
		item := NewShellItem(this, psf, pidl)
		if item.Name == "" { //?
			continue
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		hr := psf.CompareIDs(0, items[i]._relPidl, items[j]._relPidl)
		if win32.FAILED(hr) {
			return false
		}
		return int16(hr) < 0
	})

	for _, item := range items {
		win32.CoTaskMemFree(unsafe.Pointer(item._relPidl))
		item._relPidl = nil
	}

	return items
}

func (this *ShellItem) Open() {
	wc := forms.NewWaitCursor()
	defer wc.Restore()

	var sei win32.SHELLEXECUTEINFO
	sei.CbSize = uint32(unsafe.Sizeof(sei))
	sei.FMask = win32.SEE_MASK_IDLIST
	sei.Hwnd = forms.HWndActive
	sei.LpVerb = win32.StrToPwstr("open")
	sei.NShow = int32(win32.SW_SHOW)

	pidl := this.getAbsPidl()
	sei.LpIDList = unsafe.Pointer(pidl)
	ok, errno := win32.ShellExecuteEx(&sei)
	if ok != win32.TRUE {
		println(errno.Error())
	}
}

func (this *ShellItem) GetPath() string {
	wsz := make([]uint16, win32.MAX_PATH)
	ret := win32.SHGetPathFromIDList(this.getAbsPidl(), &wsz[0])
	var path string
	if ret == win32.TRUE {
		path = win32.WstrToStr(wsz)
	} else {
		pathItems := this.GetPathItems()
		path = ""
		for _, it := range pathItems {
			if len(path) > 0 {
				path += ")"
			}
			path += "(" + it.Name
		}
		path += ")"
	}
	return path
}

func (this *ShellItem) ShowContextMenu() {
	win := forms.GetActiveWin()

	pidl := this.getAbsPidl()
	var psf *win32.IShellFolder
	var pidlChild *win32.ITEMIDLIST

	hr := win32.SHBindToParent(pidl, &win32.IID_IShellFolder,
		unsafe.Pointer(&psf), &pidlChild)
	win32.ASSERT_SUCCEEDED(hr)
	defer psf.Release()

	var pcm *win32.IContextMenu
	hr = psf.GetUIObjectOf(win.GetHandle(), 1, &pidlChild,
		&win32.IID_IContextMenu, nil, unsafe.Pointer(&pcm))
	win32.ASSERT_SUCCEEDED(hr)
	defer pcm.Release()

	var pt win32.POINT
	win32.GetCursorPos(&pt)

	hMenu, _ := win32.CreatePopupMenu()
	defer win32.DestroyMenu(hMenu)
	hr = pcm.QueryContextMenu(hMenu, 0, 1, 0x6fff, win32.CMF_NORMAL)
	win32.ASSERT_SUCCEEDED(hr)

	var pcm2 *win32.IContextMenu2
	pcm.QueryInterface(&win32.IID_IContextMenu2, unsafe.Pointer(&pcm2))
	if pcm2 != nil {
		defer pcm2.Release()
	}

	var pcm3 *win32.IContextMenu3
	pcm.QueryInterface(&win32.IID_IContextMenu3, unsafe.Pointer(&pcm3))
	if pcm3 != nil {
		defer pcm3.Release()
	}

	msgProcessor := forms.MessageProcessorByFunc(func(m *forms.Message) {
		if pcm3 != nil {
			hr := pcm3.HandleMenuMsg2(m.UMsg, m.WParam, m.LParam, &m.Result)
			if win32.SUCCEEDED(hr) {
				m.Handled = true
			}
		} else if pcm2 != nil {
			hr := pcm2.HandleMenuMsg(m.UMsg, m.WParam, m.LParam)
			if win32.SUCCEEDED(hr) {
				m.Result = 0
				m.Handled = true
			}
		}
	})
	win.AddMessageProcessor(msgProcessor)
	defer win.RemoveMessageProcessor(msgProcessor)

	iCmd, errno := win32.TrackPopupMenuEx(hMenu, uint32(win32.TPM_RETURNCMD),
		pt.X, pt.Y, win.GetHandle(), nil)
	if errno != win32.NO_ERROR {
		log.Panic(errno)
	}
	//println(iCmd)
	if iCmd > 0 {
		var info win32.CMINVOKECOMMANDINFOEX
		info.CbSize = uint32(unsafe.Sizeof(info))

		const CMIC_MASK_UNICODE = 0x00004000
		info.FMask = win32.CMIC_MASK_PTINVOKE | CMIC_MASK_UNICODE
		if win32.GetKeyState(int32(win32.VK_CONTROL)) < 0 {
			info.FMask |= win32.CMIC_MASK_CONTROL_DOWN
		}
		if win32.GetKeyState(int32(win32.VK_SHIFT)) < 0 {
			info.FMask |= win32.CMIC_MASK_SHIFT_DOWN
		}
		info.Hwnd = win.GetHandle()
		info.LpVerb = win32.MAKEINTRESOURCEA(uint16(iCmd - 1))
		info.LpVerbW = win32.MAKEINTRESOURCE(uint16(iCmd - 1))
		info.NShow = int32(win32.SW_SHOWNORMAL)
		info.PtInvoke = pt
		hr = pcm.InvokeCommand((*win32.CMINVOKECOMMANDINFO)(unsafe.Pointer(&info)))
		if win32.FAILED(hr) {
			println(win32.HRESULT_ToString(hr))
		}
	}
}
