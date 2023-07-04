package forms

type Input interface {
	GetValue() any
	SetValue(any)
	Focus()
	GetOnValueChange() *SimpleEvent //todo:clearify only fire on user action?
}
