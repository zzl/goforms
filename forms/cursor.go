package forms

import (
	"bytes"
	"github.com/zzl/go-com/com"
	"github.com/zzl/go-com/com/comimpl"
	"github.com/zzl/go-win32api/v2/win32"
	"unsafe"
)

type Cursor struct {
	Handle win32.HCURSOR
	shared bool
}

func NewCursorFromData(data []byte) (*Cursor, error) {
	var pUnk *win32.IUnknown
	hr := win32.OleCreatePictureIndirect(nil, &win32.IID_IPicture, win32.TRUE, unsafe.Pointer(&pUnk))
	if !win32.SUCCEEDED(hr) {
		return nil, com.Error(hr)
	}
	pPicture := (*win32.IPicture)(unsafe.Pointer(pUnk))
	defer pPicture.Release()

	var pPersistStream *win32.IPersistStream
	hr = pUnk.QueryInterface(&win32.IID_IPersistStream, unsafe.Pointer(&pPersistStream))
	if !win32.SUCCEEDED(hr) {
		return nil, com.Error(hr)
	}
	defer pPersistStream.Release()

	pStream := comimpl.NewReaderWriterIStream(bytes.NewReader(data), nil)
	hr = pPersistStream.Load(pStream)
	if !win32.SUCCEEDED(hr) {
		return nil, com.Error(hr)
	}
	defer pStream.Release()

	var hCur uint32
	hr = pPicture.Get_Handle(&hCur)
	if !win32.SUCCEEDED(hr) {
		return nil, com.Error(hr)
	}
	hCur2, errno := win32.CopyImage(win32.HANDLE(hCur), win32.IMAGE_CURSOR, 0, 0, 0)
	if errno != win32.NO_ERROR {
		return nil, errno
	}
	return &Cursor{Handle: hCur2}, nil
}

func (this *Cursor) Dispose() {
	if this.shared || this.Handle == 0 {
		return
	}
	win32.DestroyCursor(this.Handle)
}
