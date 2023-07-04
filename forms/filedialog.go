package forms

import (
	"github.com/zzl/go-win32api/v2/win32"
	"unsafe"
)

type FileTypeFilter struct {
	Name    string
	Pattern string
}

type FileDialog struct {
	Title              string
	Filters            []FileTypeFilter
	DefaultFilterIndex int //0 based
	DefaultFileName    string
	InitialDir         string
	ResultFileName     string

	lastFilterIndex int
}

var ImagesFilter = FileTypeFilter{
	Name:    "Image files(*.gif;*.jpg;*.jpeg;*.bmp;*.png)",
	Pattern: "*.gif;*.jpg;*.jpeg;*.bmp;*.png"}

var ImagesFilters = []FileTypeFilter{ImagesFilter}

func (this *FileDialog) buildFilterStr() string {
	var str string
	for _, filter := range this.Filters {
		str += filter.Name + "\000" + filter.Pattern + "\000"
	}
	return str
}

func (this *FileDialog) buildOfn() win32.OPENFILENAME {
	var ofn win32.OPENFILENAME
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))

	buf := make([]uint16, win32.MAX_PATH)
	ofn.LpstrFile = &buf[0]
	ofn.NMaxFile = win32.MAX_PATH
	ofn.LpstrFilter = win32.StrToPwstr(this.buildFilterStr())

	filterIndex := this.DefaultFilterIndex + 1
	if this.lastFilterIndex != 0 {
		filterIndex = this.lastFilterIndex
	}
	ofn.NFilterIndex = uint32(filterIndex) //?
	if this.Title != "" {
		ofn.LpstrTitle = win32.StrToPwstr(this.Title)
	}
	if this.DefaultFileName != "" {
		copy(buf, StrToWsz(this.DefaultFileName))
	}
	ofn.Flags = win32.OFN_PATHMUSTEXIST | win32.OFN_FILEMUSTEXIST | win32.OFN_EXPLORER
	return ofn
}

type OpenFileDialog struct {
	FileDialog
	FileMustExists bool
}

func NewOpenFileDialog() *OpenFileDialog {
	return &OpenFileDialog{
		//
	}
}

func (this *OpenFileDialog) Show(hWndOwner win32.HWND) bool {
	ofn := this.buildOfn()
	ofn.HwndOwner = hWndOwner
	bOk := win32.GetOpenFileName(&ofn)
	this.lastFilterIndex = int(ofn.NFilterIndex)
	if bOk == win32.TRUE {
		this.ResultFileName = win32.PwstrToStr(ofn.LpstrFile)
		return true
	}
	return false
}

type SaveFileDialog struct {
	FileDialog
}

func NewSaveFileDialog() *SaveFileDialog {
	return &SaveFileDialog{
		//
	}
}

func (this *SaveFileDialog) Show() bool {
	ofn := this.buildOfn()
	bOk := win32.GetSaveFileName(&ofn)
	this.lastFilterIndex = int(ofn.NFilterIndex)
	if bOk == win32.TRUE {
		this.ResultFileName = win32.PwstrToStr(ofn.LpstrFile)
		return true
	}
	return false
}
