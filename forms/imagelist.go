package forms

import (
	"errors"
	"io/ioutil"
	"log"
	"unsafe"

	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/drawing"
)

type ImageList struct {
	Cx     int
	Cy     int
	Bpp    byte
	Masked bool

	handle win32.HIMAGELIST
	owned  bool
}

func (this *ImageList) Init() {
}

func (this *ImageList) Dispose() {
	if this.handle == 0 {
		return
	}
	if this.owned {
		win32.ImageList_Destroy(this.handle)
	}
	this.handle = 0
}

func NewImageList(size int, create bool) *ImageList {
	iml := &ImageList{
		Cx:     size,
		Cy:     size,
		Bpp:    32,
		Masked: false,
	}
	iml.Init()
	if create {
		err := iml.Create()
		if err != nil {
			log.Panic(err)
		}
	}
	return iml
}

func NewImageListEx(size int, bpp int, masked bool) *ImageList {
	iml := &ImageList{
		Cx:     size,
		Cy:     size,
		Bpp:    byte(bpp),
		Masked: masked,
	}
	iml.Init()
	return iml
}

func (this *ImageList) CopyUnowned() *ImageList {
	iml := *this
	iml.owned = false
	return &iml
}

func GetSysImageList(small bool) *ImageList {
	var size = 32
	iImageList := int32(win32.SHIL_LARGE)
	if small {
		iImageList = int32(win32.SHIL_SMALL)
		size = 16
	}
	var hIml win32.HIMAGELIST
	hr := win32.SHGetImageList(iImageList, &win32.IID_IImageList, unsafe.Pointer(&hIml))
	if win32.FAILED(hr) {
		log.Fatal("?")
	}

	return &ImageList{
		handle: hIml,
		Cx:     size,
		Cy:     size,
	}
}

func NewImageListFromHandle(hIml win32.HIMAGELIST, owned bool) *ImageList {
	var cx, cy int32
	ret := win32.ImageList_GetIconSize(hIml, &cx, &cy)
	if ret == win32.FALSE {
		log.Println("?")
	}
	iml := &ImageList{
		Cx:     int(cx),
		Cy:     int(cy),
		Bpp:    0, //?
		Masked: false,
		handle: hIml,
	}
	iml.Init()
	iml.owned = owned
	return iml
}

func CreateImageListFromIcons(size int, hIcons ...win32.HICON) (*ImageList, error) {
	iml := &ImageList{
		Cx:     size,
		Cy:     size,
		Bpp:    32,
		Masked: false,
	}
	iml.Init()
	err := iml.Create()
	if err != nil {
		return nil, err
	}
	for _, hIcon := range hIcons {
		iml.AddIcon(hIcon)
	}
	return iml, nil
}

func CreateImageListFromImage(imageData []byte) (*ImageList, error) {
	bitmap, _ := drawing.NewBitmapFromBytes(nil, imageData)
	cx, cy := bitmap.GetSizeDecomposed()
	size := min(int(cx), int(cy))

	iml := &ImageList{
		Cx:     size,
		Cy:     size,
		Bpp:    32,
		Masked: false,
	}
	iml.Init()
	err := iml.Create()
	if err != nil {
		return nil, err
	}
	hBitmap := bitmap.GetHBitmap()
	_, err = iml.Add(hBitmap)
	win32.DeleteObject(win32.HGDIOBJ(hBitmap))
	if err != nil {
		iml.Dispose()
		return nil, err
	}
	return iml, nil
}

func (this *ImageList) Create() error {
	if this.Cx == 0 {
		this.Cx = 16
	}
	if this.Cy == 0 {
		this.Cy = this.Cx
	}
	if this.Bpp == 0 {
		this.Bpp = 32
	}
	flags := uint32(this.Bpp)
	if this.Masked {
		flags |= uint32(win32.ILC_MASK)
	}
	hIml := win32.ImageList_Create(
		int32(this.Cx), int32(this.Cy), win32.IMAGELIST_CREATION_FLAGS(flags), 4, 4)
	if hIml == 0 {
		return errors.New("?")
	}
	this.handle = hIml
	this.Init()
	this.owned = true
	return nil
}

func (this *ImageList) GetHandle() win32.HIMAGELIST {
	return this.handle
}

func (this *ImageList) Detach() win32.HIMAGELIST {
	handle := this.handle
	this.handle = 0
	return handle
}

func (this *ImageList) AddImageDatas(imageData ...[]byte) (int, error) {
	var ret int
	var err error
	for _, data := range imageData {
		ret, err = this.AddImageData(data)
	}
	return ret, err
}

func (this *ImageList) AddImageData(imageData []byte) (int, error) {
	bitmap, _ := drawing.NewBitmapFromBytes(nil, imageData)
	defer bitmap.Dispose()

	hBitmap := bitmap.GetHBitmap()
	ret, errno := this.Add(hBitmap)
	win32.DeleteObject(hBitmap)

	return ret, errno
}

func (this *ImageList) AddImageFile(filePath string) (int, error) {
	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		return -1, err
	}
	return this.AddImageData(bts)
}

func (this *ImageList) Add(hBitmap win32.HBITMAP) (int, error) {
	ret := win32.ImageList_Add(this.handle, hBitmap, 0)
	if ret == -1 {
		return -1, errors.New("?")
	}
	return int(ret), nil
}

func (this *ImageList) AddIcon(hIcon win32.HICON) (int, error) {
	ret, errno := win32.ImageList_AddIcon(this.handle, hIcon)
	return int(ret), errno.NilOrError()
}

func (this *ImageList) GetCount() int {
	count := win32.ImageList_GetImageCount(this.handle)
	return int(count)
}

func (this *ImageList) Resize(cx int) *ImageList {
	iml2 := &ImageList{
		Cx:     cx,
		Cy:     this.Cy,
		Bpp:    32,
		Masked: false,
	}
	iml2.Init()
	iml2.Create()
	count := this.GetCount()

	hdcScreen := win32.GetDC(0)
	hdc := win32.CreateCompatibleDC(hdcScreen)
	var rc win32.RECT
	rc.Right = int32(cx)
	rc.Bottom = int32(this.Cy)
	hbr := win32.GetSysColorBrush(win32.COLOR_WINDOW)
	for n := 0; n < count; n++ {
		//iml2
		hbmp := win32.CreateCompatibleBitmap(hdcScreen, int32(cx), int32(this.Cy))
		hOriBmp := win32.SelectObject(hdc, win32.HGDIOBJ(hbmp))
		win32.FillRect(hdc, &rc, hbr)

		win32.ImageList_DrawEx(this.handle, int32(n), hdc,
			int32(cx-this.Cx), 0, 0, 0,
			win32.CLR_NONE_U, win32.CLR_NONE_U, win32.ILD_NORMAL)

		win32.SelectObject(hdc, hOriBmp)
		win32.ImageList_Add(iml2.handle, hbmp, 0)
		win32.DeleteObject(win32.HGDIOBJ(hbmp))

	}
	win32.DeleteDC(hdc)
	win32.ReleaseDC(0, hdcScreen)
	return iml2
}

func (this *ImageList) GetHBitmap(index int) win32.HBITMAP {
	hIcon := win32.ImageList_GetIcon(this.handle, int32(index), uint32(win32.ILD_NORMAL))
	if hIcon == 0 {
		log.Println("?")
	}
	bitmap := drawing.NewBitmapFromHIcon(nil, hIcon)
	hBitmap := bitmap.GetHBitmap()
	bitmap.Dispose()
	return hBitmap
}

func (this *ImageList) GetHIcon(index int) win32.HICON {
	hIcon := win32.ImageList_GetIcon(this.handle, int32(index), uint32(win32.ILD_NORMAL))
	return hIcon
}

func (this *ImageList) Draw(index int, hdc win32.HDC, x int32, y int32) {
	win32.ImageList_DrawEx(this.handle, int32(index), hdc,
		x, y, int32(this.Cx), int32(this.Cy),
		win32.CLR_NONE_U, win32.CLR_NONE_U, win32.ILD_NORMAL)
}
