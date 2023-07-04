package drawing

import (
	"github.com/zzl/go-gdiplus/gdip"
	"github.com/zzl/go-win32api/v2/win32"
	"runtime"
)

type StringFormat struct {
	p *gdip.StringFormat
}

func newStringFormat(s *Scope, p *gdip.StringFormat) *StringFormat {
	sf := &StringFormat{p}
	if s != nil {
		s.Add(sf)
	}
	runtime.SetFinalizer(sf, (*StringFormat).Dispose)
	return sf
}

func NewStringFormat(s *Scope) *StringFormat {
	return NewStringFormatOptsLang(s, 0, 0)
}

func NewStringFormatOpts(s *Scope, options gdip.StringFormatFlags) *StringFormat {
	return NewStringFormatOptsLang(s, options, 0)
}

func NewStringFormatOptsLang(s *Scope, options gdip.StringFormatFlags,
	language int32) *StringFormat {
	var pFormat *gdip.StringFormat
	status := gdip.CreateStringFormat(int32(options), win32.LANGID(language), &pFormat)
	checkStatus(status)
	return newStringFormat(s, pFormat)
}

func GetGenericDefaultStringFormat(s *Scope) *StringFormat {
	var pFormat *gdip.StringFormat
	status := gdip.StringFormatGetGenericDefault(&pFormat)
	checkStatus(status)
	return newStringFormat(s, pFormat)
}

func GetGenericTypographicStringFormat(s *Scope) *StringFormat {
	var pFormat *gdip.StringFormat
	status := gdip.StringFormatGetGenericTypographic(&pFormat)
	checkStatus(status)
	return newStringFormat(s, pFormat)
}

func (this *StringFormat) Clone(s *Scope) *StringFormat {
	var pFormat2 *gdip.StringFormat
	status := gdip.CloneStringFormat(this.p, &pFormat2)
	checkStatus(status)
	return newStringFormat(s, pFormat2)
}

func (this *StringFormat) Dispose() {
	if this.p == nil {
		return
	}
	status := gdip.DeleteStringFormat(this.p)
	checkStatus(status)
	this.p = nil
	runtime.SetFinalizer(this, nil)
}

func (this *StringFormat) GetFormatFlags() gdip.StringFormatFlags {
	var flags int32
	status := gdip.GetStringFormatFlags(this.p, &flags)
	checkStatus(status)
	return gdip.StringFormatFlags(flags)
}

func (this *StringFormat) SetFormatFlags(flags gdip.StringFormatFlags) {
	status := gdip.SetStringFormatFlags(this.p, int32(flags))
	checkStatus(status)
}

func (this *StringFormat) SetMeasurableCharacterRanges(ranges []gdip.CharacterRange) {
	status := gdip.SetStringFormatMeasurableCharacterRanges(
		this.p, int32(len(ranges)), &ranges[0])
	checkStatus(status)
}

func (this *StringFormat) GetAlignment() gdip.StringAlignment {
	var align gdip.StringAlignment
	status := gdip.GetStringFormatAlign(this.p, &align)
	checkStatus(status)
	return align
}

func (this *StringFormat) SetAlignment(align gdip.StringAlignment) {
	status := gdip.SetStringFormatAlign(this.p, align)
	checkStatus(status)
}

func (this *StringFormat) GetLineAlignment() gdip.StringAlignment {
	var align gdip.StringAlignment
	status := gdip.GetStringFormatLineAlign(this.p, &align)
	checkStatus(status)
	return align
}

func (this *StringFormat) SetLineAlignment(align gdip.StringAlignment) {
	status := gdip.SetStringFormatLineAlign(this.p, align)
	checkStatus(status)
}

func (this *StringFormat) GetHotkeyPrefix() gdip.HotkeyPrefix {
	var prefix int32
	status := gdip.GetStringFormatHotkeyPrefix(this.p, &prefix)
	checkStatus(status)
	return gdip.HotkeyPrefix(prefix)
}

func (this *StringFormat) SetHotkeyPrefix(prefix gdip.HotkeyPrefix) {
	status := gdip.SetStringFormatHotkeyPrefix(this.p, int32(prefix))
	checkStatus(status)
}

func (this *StringFormat) GetTabStops() (firstTabOffset float32, tabStops []float32) {
	var count int32
	status := gdip.GetStringFormatTabStopCount(this.p, &count)
	checkStatus(status)
	tabStops = make([]float32, count)
	status = gdip.GetStringFormatTabStops(this.p, count, &firstTabOffset, &tabStops[0])
	checkStatus(status)
	return
}

func (this *StringFormat) SetTabStops(firstTabOffset float32, tabStops []float32) {
	status := gdip.SetStringFormatTabStops(this.p,
		firstTabOffset, int32(len(tabStops)), &tabStops[0])
	checkStatus(status)
}

func (this *StringFormat) GetTrimming() gdip.StringTrimming {
	var trimming gdip.StringTrimming
	status := gdip.GetStringFormatTrimming(this.p, &trimming)
	checkStatus(status)
	return trimming
}

func (this *StringFormat) SetTrimming(trimming gdip.StringTrimming) {
	status := gdip.SetStringFormatTrimming(this.p, trimming)
	checkStatus(status)
}

func (this *StringFormat) SetDigitSubstitution(language win32.LANGID,
	substitute gdip.StringDigitSubstitute) {
	status := gdip.SetStringFormatDigitSubstitution(this.p,
		language, substitute)
	checkStatus(status)
}

func (this *StringFormat) GetDigitSubstitution() (win32.LANGID, gdip.StringDigitSubstitute) {
	var lang win32.LANGID
	var substitute gdip.StringDigitSubstitute
	status := gdip.GetStringFormatDigitSubstitution(this.p, &lang, &substitute)
	checkStatus(status)
	return lang, substitute
}
