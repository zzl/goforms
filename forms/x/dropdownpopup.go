package forms

import . "github.com/zzl/goforms/forms"

type DropdownPopup interface {
	GetControl() Control

	GetOnOk() *SimpleEvent
	GetOnCancel() *SimpleEvent

	GetValue() interface{}
	SetValue(value interface{})

	GetText() string

	PreparePopup()

	GetPopupSize(width int, maxWidth int, maxHeight int) (int, int)

	NotifyBeforeShow()
	NotifyAfterShow()

	SetContainer(container DropdownPopupContainer)

	//HandleNotify(info *NotifyInfo)

}
