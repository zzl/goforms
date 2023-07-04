package events

import "reflect"

type EventInfo interface {
	GetSender() any
	SetSender(sender any)
	GetHandled() bool
	SetHandled(handled bool)
	GetResult() uintptr
}

type EventListener[T EventInfo] func(ei T)
type Event[T EventInfo] struct {
	listeners []*EventListener[T]
}

type SimpleEventListener = EventListener[*SimpleEventInfo]

func (this *Event[T]) AddListener(listener EventListener[T]) *EventListener[T] {
	p := &listener
	this.listeners = append(this.listeners, p)
	return p
}

func (this *Event[T]) RemoveListener(pListener *EventListener[T]) {
	for n, listener := range this.listeners {
		if listener == pListener {
			this.listeners = append(this.listeners[:n], this.listeners[n+1:]...)
			return
		}
	}
}

func (this *Event[T]) Fire(sender any, eventInfo T) {
	if this.listeners == nil {
		return
	}
	if !reflect.ValueOf(eventInfo).IsNil() {
		eventInfo.SetSender(sender)
	}
	for _, listener := range this.listeners {
		(*listener)(eventInfo)
	}
}

type SimpleEventInfo struct {
	Sender  any
	Handled bool
	Result  uintptr
}

func (this *SimpleEventInfo) GetSender() any {
	return this.Sender
}

func (this *SimpleEventInfo) SetSender(sender any) {
	if this == nil { //..
		return
	}
	this.Sender = sender
}

func (this *SimpleEventInfo) GetHandled() bool {
	return this.Handled
}

func (this *SimpleEventInfo) SetHandled(handled bool) {
	this.Handled = handled
}

func (this *SimpleEventInfo) GetResult() uintptr {
	return this.Result
}

type SimpleEvent = Event[*SimpleEventInfo]

type ExtraEventInfo struct {
	SimpleEventInfo
	Extra map[string]interface{}
}

func NewExtraEventInfo(key string, value interface{}) *ExtraEventInfo {
	return &ExtraEventInfo{
		Extra: map[string]interface{}{
			key: value,
		},
	}
}

type ExtraEvent = Event[*ExtraEventInfo]
