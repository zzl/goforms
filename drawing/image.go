package drawing

import (
	"errors"
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

type Image struct {
	p *gdip.Image
}

func newImage(s *Scope, p *gdip.Image) *Image {
	image := &Image{p}
	if s != nil {
		s.Add(image)
	}
	runtime.SetFinalizer(image, (*Image).Dispose)
	return image
}

func NewImage(p *gdip.Image) *Image {
	return &Image{p}
}

func NewImageFromFile(s *Scope, filename string) (*Image, error) {
	pwsz, _ := syscall.UTF16PtrFromString(filename)
	var pImage *gdip.Image
	status := gdip.LoadImageFromFile(pwsz, &pImage)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newImage(s, pImage), nil
}

func NewImageFromBytes(s *Scope, bytes []byte) (*Image, error) {
	pStream := win32.SHCreateMemStream(&bytes[0], uint32(len(bytes)))
	var pImage *gdip.Image
	status := gdip.LoadImageFromStream(pStream, &pImage)
	pStream.Release()
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newImage(s, pImage), nil
}

func NewImageFromHBitmap(s *Scope, hBitmap win32.HBITMAP) *Image {
	if hBitmap == 0 {
		return nil
	}
	bitmap := NewBitmapFromHBitmap(s, hBitmap)
	return bitmap.AsImage()
}

func (this *Image) P() *gdip.Image {
	return this.p
}

func (this *Image) Clone(s *Scope) *Image {
	var pImage *gdip.Image
	status := gdip.CloneImage(this.p, &pImage)
	checkStatus(status)
	return newImage(s, pImage)
}

func (this *Image) Dispose() {
	if this.p == nil {
		return
	}
	gdip.DisposeImage(this.p)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *Image) IsValid() bool {
	return this.p != nil
}

func (this *Image) GetHBitmap() win32.HBITMAP {
	var hbm win32.HBITMAP
	status := gdip.CreateHBITMAPFromBitmap((*gdip.Bitmap)(unsafe.Pointer(this.p)), &hbm, 0)
	checkStatus(status)
	return hbm
}

func (this *Image) Save(filename string) error {
	return this.SaveWithParams(filename, nil)
}

func (this *Image) SaveWithParams(filename string, params *gdip.EncoderParameters) error {
	ext := strings.ToLower(filepath.Ext(filename))
	clsid, ok := extClsIdMap[ext]
	if !ok {
		return errors.New("unsupported format")
	}
	pwsz, _ := syscall.UTF16PtrFromString(filename)
	status := gdip.SaveImageToFile(this.p, pwsz, &clsid, params)
	if status != gdip.Ok {
		return GdipError(status)
	}
	return nil
}

func (this *Image) WriteTo(writer io.Writer, formatExt string) error {
	return this.WriteToWithParams(writer, formatExt, nil)
}

func (this *Image) WriteToWithParams(writer io.Writer, formatExt string,
	params *gdip.EncoderParameters) error {
	ext := strings.ToLower(formatExt)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	clsid, ok := extClsIdMap[ext]
	if !ok {
		return errors.New("unsupported format")
	}

	_ = clsid
	//xxxx
	//stream := com.NewStream(nil, writer)
	//status := gdip.SaveImageToStream(this.p, stream.ComInterface, &clsid, params)
	//if status != gdip.Ok {
	//	return GdipError(status)
	//}
	return nil
}

func (this *Image) GetSize() Size {
	var w, h float32
	status := gdip.GetImageDimension(this.p, &w, &h)
	ensureOk(status)
	return Size{int32(w), int32(h)}
}

func (this *Image) GetSizeDecomposed() (width, height int32) {
	size := this.GetSize()
	return size.Width, size.Height
}

func (this *Image) GetPixelFormat() gdip.PixelFormat {
	var format gdip.PixelFormat
	status := gdip.GetImagePixelFormat(this.p, &format)
	checkStatus(status)
	return format
}

func (this *Image) GetThumbnail(s *Scope, thumbWidth, thumbHeight uint32) (*Image, error) {
	var pImage *gdip.Image
	status := gdip.GetImageThumbnail(this.p, thumbWidth, thumbHeight, &pImage, 0, nil)
	if status != gdip.Ok {
		return nil, GdipError(status)
	}
	return newImage(s, pImage), nil
}

func (this *Image) RotateFlip(typ gdip.RotateFlipType) {
	status := gdip.ImageRotateFlip(this.p, typ)
	checkStatus(status)
}

func GetPixelFormatSize(format gdip.PixelFormat) int {
	return int(format>>8) & 0xFF
}

// ImageAttributes
type ImageAttributes struct {
	p *gdip.ImageAttributes
}

func newImageAttributes(s *Scope, p *gdip.ImageAttributes) *ImageAttributes {
	attrs := &ImageAttributes{p}
	if s != nil {
		s.Add(attrs)
	}
	runtime.SetFinalizer(attrs, (*ImageAttributes).Dispose)
	return attrs
}

func NewImageAttributes(s *Scope) *ImageAttributes {
	var p *gdip.ImageAttributes
	status := gdip.CreateImageAttributes(&p)
	checkStatus(status)
	return newImageAttributes(s, p)
}

func (this *ImageAttributes) Dispose() {
	if this.p == nil {
		return
	}
	gdip.DisposeImageAttributes(this.p)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *ImageAttributes) SetColorKey(colorLow Color, colorHigh Color) {
	status := gdip.SetImageAttributesColorKeys(this.p,
		gdip.ColorAdjustTypeDefault,
		win32.TRUE, colorLow.Argb(), colorHigh.Argb())
	checkStatus(status)
}

var extClsIdMap map[string]win32.CLSID

func init() {
	extClsIdMap = map[string]win32.CLSID{
		".bmp": gdip.BmpEncoderId,
		".jpg": gdip.JpgEncoderId,
		".png": gdip.PngEncoderId,
	}
}
