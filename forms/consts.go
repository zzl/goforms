package forms

import (
	"errors"
	"github.com/zzl/goforms/framework/consts"
	"math"

	"github.com/zzl/go-win32api/v2/win32"
)

const WM_REFLECT_COMMAND = win32.WM_APP + 64
const WM_REFLECT_NOTIFY = win32.WM_APP + 65
const WM_REFLECT_MEASUREITEM = win32.WM_APP + 66
const WM_REFLECT_DRAWITEM = win32.WM_APP + 67

const WM_APP_DISPATCH = win32.WM_APP + 88

const WM_CHILD_SETFOCUS = win32.WM_APP + 100
const WM_CHILD_KILLFOCUS = win32.WM_APP + 101

const NegativeOne = consts.NegativeOne
const NegativeOne32 = consts.NegativeOne32

const TransparentColor win32.COLORREF = math.MaxUint32 - 1

const NullColor win32.COLORREF = math.MaxUint32 - 2

const (
	Data_ContextMenu = "__ContextMenu"
	Data_Font        = "__Font"
	Data_FontHandle  = "__FontHandle"
	Data_Disposables = "__Disposables"
	Data_ModalResult = "__ModalResult"

	Data_BackColor      = "__BackColor"
	Data_BackColorBrush = "__BackColorBrush"

	Data_Collapsed = "__Collapsed"
)

var CancelError error = errors.New("Canceled")
