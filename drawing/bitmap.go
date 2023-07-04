package drawing

import (
	"errors"
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"github.com/zzl/goforms/framework/scope"
	"log"
	"runtime"
	"syscall"
	"unsafe"
)

type Bitmap struct {
	Image
	p *gdip.Bitmap
}

func (this *Bitmap) Dispose() {
	if this.p == nil {
		return
	}
	this.Image.Dispose()
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Bitmap) AsImage() *Image {
	return &this.Image
}

func newBitmap(s *Scope, p *gdip.Bitmap) *Bitmap {
	bitmap := &Bitmap{
		Image: Image{
			p: &p.Image,
		},
		p: p,
	}
	if s != nil {
		s.Add(bitmap)
	}
	runtime.SetFinalizer(bitmap, (*Bitmap).Dispose)
	return bitmap
}

func NewBitmapFromFile(s *Scope, filename string) (*Bitmap, error) {
	pwsz, _ := syscall.UTF16PtrFromString(filename)
	var pBitmap *gdip.Bitmap
	status := gdip.CreateBitmapFromFile(pwsz, &pBitmap)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newBitmap(s, pBitmap), nil
}

func NewBitmapFromBytes(s *Scope, bytes []byte) (*Bitmap, error) {
	pStream := win32.SHCreateMemStream(&bytes[0], uint32(len(bytes)))
	var pBitmap *gdip.Bitmap
	status := gdip.CreateBitmapFromStream(pStream, &pBitmap)
	pStream.Release()
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newBitmap(s, pBitmap), nil
}

func NewBitmapFromHIcon(s *Scope, hIcon win32.HICON) *Bitmap {
	return hiconToBitmap(s, hIcon)
}

func NewBitmapFromHBitmap(s *Scope, hBitmap win32.HBITMAP) *Bitmap {
	var pBitmap *gdip.Bitmap
	status := gdip.CreateBitmapFromHBITMAP(hBitmap, 0, &pBitmap)
	ensureOk(status)
	return newBitmap(s, pBitmap)
}

func NewBitmap(s *Scope, width int32, height int32, format gdip.PixelFormat, scan0 *byte) *Bitmap {
	var pBitmap *gdip.Bitmap
	status := gdip.CreateBitmapFromScan0(width, height, 0, format, nil, &pBitmap)
	ensureOk(status)
	return newBitmap(s, pBitmap)
}

func NewBitmapFromGraphics(s *Scope, width, height int32, g *Graphics) *Bitmap {
	var pBitmap *gdip.Bitmap
	status := gdip.CreateBitmapFromGraphics(width, height, g.p, &pBitmap)
	ensureOk(status)
	return newBitmap(s, pBitmap)
}

func NewBitmapFromImage(s *Scope, image *Image) *Bitmap {
	width, height := image.GetSizeDecomposed()
	return NewBitmapFromImageWithSize(s, image, width, height)
}

func NewBitmapFromImageWithSize(s *Scope, image *Image, width, height int32) *Bitmap {
	bitmap := NewBitmap(s, width, height, gdip.PixelFormat32bppARGB, nil)
	g, err := NewGraphicsFromImage(nil, bitmap.AsImage())
	if err != nil {
		log.Fatalln(err)
	}
	defer g.Dispose()
	g.Clear(Color{})
	g.DrawImageRect(image, 0, 0, width, height)
	return bitmap
}

func (this *Bitmap) setP(p *gdip.Bitmap) {
	this.p = p
	this.Image.p = &p.Image
}

func (this *Bitmap) GetHBitmap() win32.HBITMAP {
	var hBitmap win32.HBITMAP
	gdip.CreateHBITMAPFromBitmap(this.p, &hBitmap, 0)
	return hBitmap
}

func (this *Bitmap) GetHIcon() win32.HICON {
	var hIcon win32.HICON
	gdip.CreateHICONFromBitmap(this.p, &hIcon)
	return hIcon
}

func (this *Bitmap) CloneRegion(s *Scope, x int32, y int32, w int32, h int32) *Bitmap {
	size := this.GetSize()
	cx, cy := size.Width, size.Height
	if x < 0 {
		x = 0
	} else if x > cx {
		x = cx
	}
	if y < 0 {
		y = 0
	} else if y > cy {
		y = cy
	}
	var pf gdip.PixelFormat
	gdip.GetImagePixelFormat(this.Image.p, &pf)
	var pClone *gdip.Bitmap
	gdip.CloneBitmapAreaI(x, y,
		w, h, pf, this.p, &pClone)
	return newBitmap(s, pClone)
}

func (this *Bitmap) MakeTransparent() {
	transparent := ColorOf(0xFFFF00FF) //Colors.Magenta
	w, h := this.GetSizeDecomposed()
	if w > 0 && h > 0 {
		transparent = this.GetPixel(0, h-1)
	}
	if transparent.A() < 255 {
		return
	}
	this.MakeTransparentOnColor(transparent)
}

func (this *Bitmap) MakeTransparentOnColor(color Color) {
	s := scope.NewScope() //leave?
	w, h := this.GetSizeDecomposed()
	result := NewBitmap(s, w, h, gdip.PixelFormat32bppARGB, nil)

	g, _ := NewGraphicsFromBitmap(s, result)
	g.Clear(Color{})
	rect := Rect{0, 0, w, h}

	a := NewImageAttributes(s)
	a.SetColorKey(color, color)
	g.DrawImageRectRectAttr(this.AsImage(), rect, rect, gdip.UnitPixel, a)

	p := this.p
	this.setP(result.p)
	result.setP(p)
}

func (this *Bitmap) LockBits(rect Rect, flags gdip.ImageLockMode,
	format gdip.PixelFormat) gdip.BitmapData {
	var data gdip.BitmapData
	status := gdip.BitmapLockBits(this.p, (*gdip.Rect)(&rect),
		uint32(flags), format, &data)
	checkStatus(status)
	return data
}

func (this *Bitmap) UnlockBits(data gdip.BitmapData) {
	status := gdip.BitmapUnlockBits(this.p, &data)
	checkStatus(status)
}

func (this *Bitmap) GetPixel(x int32, y int32) Color {
	var pix gdip.ARGB
	status := gdip.BitmapGetPixel(this.p, x, y, &pix)
	checkStatus(status)
	return ColorOf(pix)
}

func (this *Bitmap) SetPixel(x int32, y int32, color Color) {
	status := gdip.BitmapSetPixel(this.p, x, y, color.Argb())
	checkStatus(status)
}

///

func (this *Bitmap) SaveToFile(filePath string) error {
	pwsz, _ := syscall.UTF16PtrFromString(filePath)

	//var sGuid string
	//dotPos := strings.LastIndexByte(filePath, '.')
	//ext := filePath[dotPos:]
	//if ext == ".jpg" {
	//	sGuid = "{557CF401-1A04-11D3-9A73-0000F81EF32E}"
	//} else if ext == ".png" {
	//	sGuid = "{557cf406-1a04-11d3-9a73-0000f81ef32e}"
	//} else if ext == ".bmp" {
	//	sGuid = "{557cf400-1a04-11d3-9a73-0000f81ef32e}"
	//} else {
	//	return errors.New("unsupported format")
	//}
	//guid, _ := windows.GUIDFromString(sGuid)
	//clsid := (*win32.CLSID)(&guid)
	clsid := &gdip.BmpEncoderId
	status := gdip.SaveImageToFile(&this.p.Image, pwsz, clsid, nil)
	if status != gdip.Ok {
		return errors.New("??")
	} else {
		return nil
	}
}

// ?
func hiconToBitmap(s *Scope, hIcon win32.HICON) *Bitmap {
	var ii win32.ICONINFO
	win32.GetIconInfo(hIcon, &ii)

	var bmp win32.BITMAP
	win32.GetObject(win32.HGDIOBJ(ii.HbmColor), int32(unsafe.Sizeof(bmp)), unsafe.Pointer(&bmp))

	win32.DeleteObject(win32.HGDIOBJ(ii.HbmMask))
	win32.DeleteObject(win32.HGDIOBJ(ii.HbmColor))

	width, height := bmp.BmWidth, bmp.BmHeight

	//
	var format gdip.PixelFormat = 0x26200A //Format32bppArgb
	//var format gdip.PixelFormat = 0x000e200b
	var pBitmap *gdip.Bitmap
	status := gdip.CreateBitmapFromScan0(width, height, 0, format, nil, &pBitmap)
	if status != gdip.Ok {
		log.Fatal(status)
	}

	var pGraphics *gdip.Graphics
	status = gdip.GetImageGraphicsContext(&pBitmap.Image, &pGraphics)
	if status != gdip.Ok {
		log.Fatal(status)
	}

	gdip.GraphicsClear(pGraphics, 0x00FFFFFF)
	var hdc win32.HDC
	gdip.GetDC(pGraphics, &hdc)

	ret, errno := win32.DrawIconEx(win32.HDC(hdc),
		0, 0, hIcon, width, height, 0, 0, win32.DI_NORMAL)

	if ret == 0 {
		log.Println(errno.Error())
	}
	gdip.ReleaseDC(pGraphics, hdc)

	return newBitmap(s, pBitmap)
}

func HIconToHBitmap(hIcon win32.HICON) win32.HBITMAP {
	var ii win32.ICONINFO
	bOk, errno := win32.GetIconInfo(hIcon, &ii)
	if bOk == 0 {
		log.Println(errno)
		return 0
	}
	defer func() {
		if ii.HbmColor != 0 {
			win32.DeleteObject(win32.HGDIOBJ(ii.HbmColor))
		}
		if ii.HbmMask != 0 {
			win32.DeleteObject(win32.HGDIOBJ(ii.HbmMask))
		}
	}()
	var bi win32.BITMAPINFO
	bh := &bi.BmiHeader
	bh.BiSize = uint32(unsafe.Sizeof(bi))
	bh.BiPlanes = 1
	bh.BiBitCount = 32
	bh.BiCompression = win32.BI_RGB

	if ii.HbmColor != 0 {
		var bm win32.BITMAP
		win32.GetObject(win32.HGDIOBJ(ii.HbmColor),
			int32(unsafe.Sizeof(bm)), unsafe.Pointer(&bm))
		bh.BiWidth = bm.BmWidth
		bh.BiHeight = bm.BmHeight
	} else { //?
		var bm win32.BITMAP
		win32.GetObject(win32.HGDIOBJ(ii.HbmMask),
			int32(unsafe.Sizeof(bm)), unsafe.Pointer(&bm))
		bh.BiWidth = bm.BmWidth
		bh.BiHeight = bm.BmHeight / 2
	}
	if bh.BiWidth == 0 {
		return 0
	}
	hdc := win32.CreateCompatibleDC(0)
	var pBits unsafe.Pointer
	dib, errno := win32.CreateDIBSection(hdc, &bi,
		win32.DIB_RGB_COLORS, unsafe.Pointer(&pBits), 0, 0)
	if dib == 0 {
		win32.DeleteDC(hdc)
		return 0
	}
	hOriBitmap := win32.SelectObject(hdc, win32.HGDIOBJ(dib))

	rc := win32.RECT{0, 0, bh.BiWidth, bh.BiHeight}
	//hbrBg, _ := win32.GetStockObject(win32.WHITE_BRUSH)
	//win32.FillRect(hdc, &rc, win32.HBRUSH(hbrBg))
	ret, errno := win32.DrawIconEx(hdc,
		0, 0, hIcon, rc.Right, rc.Bottom, 0, 0, win32.DI_NORMAL)
	if ret == 0 {
		log.Println(errno) //?
	}
	//if mask?
	win32.SelectObject(hdc, hOriBitmap)
	win32.DeleteDC(hdc)
	//
	return dib
}

// ?
// upsidedown..
func HBitmapToBitmap(hBitmap win32.HBITMAP) *Bitmap {

	var bm win32.BITMAP
	ret := win32.GetObject(win32.HGDIOBJ(hBitmap),
		int32(unsafe.Sizeof(bm)), unsafe.Pointer(&bm))
	if ret == 0 {
		log.Println("??")
		return nil
	}
	var format gdip.PixelFormat = 0x26200A //Format32bppArgb
	var pBitmap *gdip.Bitmap
	status := gdip.CreateBitmapFromScan0(
		bm.BmWidth, bm.BmHeight, 0, format, nil, &pBitmap)
	if status != gdip.Ok {
		log.Println(status)
		return nil
	}
	rc := Rect{0, 0, bm.BmHeight, bm.BmHeight}
	const ImageLockModeWrite = 0x0002
	var bitmapData gdip.BitmapData
	status = gdip.BitmapLockBits(pBitmap, (*gdip.Rect)(&rc),
		ImageLockModeWrite, format, &bitmapData)
	if status != gdip.Ok {
		//?
		log.Println(status)
		return nil
	}
	if bitmapData.Stride != bm.BmWidthBytes {
		log.Println("wrong pixel format..")
		return nil
	}
	pScan0 := unsafe.Pointer(bitmapData.Scan0)
	win32.CopyMemory(pScan0, bm.BmBits,
		uint32(bm.BmWidthBytes*bm.BmHeight))

	status = gdip.BitmapUnlockBits(pBitmap, &bitmapData)
	if status != gdip.Ok {
		//?
		log.Println(status)
		return nil
	}
	return newBitmap(nil, pBitmap)
}
