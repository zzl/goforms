package layouts

import (
	"github.com/zzl/goforms/framework/events"
	"github.com/zzl/goforms/layouts/aligns"
)

const (
	Data_Layout = "__Layout"

	DefaultItemAlign    = aligns.Center
	DefaultContentAlign = aligns.Stretch
)

type LayoutAware interface {
	SetLayout(layout Layout)
}

type Collapsible interface {
	SetCollapsed(collapsed bool)
	IsCollapsed() bool
}

type CollapsibleObject struct {
	Collapsed bool
}

func (this *CollapsibleObject) SetCollapsed(collapsed bool) {
	this.Collapsed = collapsed
}

func (this *CollapsibleObject) IsCollapsed() bool {
	return this.Collapsed
}

type LayoutItem interface {
	Collapsible

	GetControl() Control
	GetName() string
	GetLayout() Layout
	SetWidth(value int)  //?
	SetHeight(value int) //?
}

type LayoutEventSource interface {
	GetOnPreLayout() *LayoutEvent
	GetOnPostLayout() *LayoutEvent
}

type LayoutEventInfo struct {
	SimpleEventInfo
	Bounds Rect
}

type LayoutEvent = events.Event[*LayoutEventInfo]

type Layout interface {
	BoundsAware
	LayoutEventSource

	SetContainer(container Container)
	SetItemDefaults(itemDefaults LayoutItem)
	AddItems(items []LayoutItem, prepend bool)
	Clone() Layout

	GetPreferredSize(layoutWidth int, layoutHeight int) (int, int)
	FindItemByControl(control Control) LayoutItem
	GetItem(name string) LayoutItem

	SetSizeGroup(map[string]int)
	Update()
}

type BaseLayout struct {
	OnPreLayout  LayoutEvent
	OnPostLayout LayoutEvent
}

func (this *BaseLayout) GetOnPreLayout() *LayoutEvent {
	return &this.OnPreLayout
}

func (this *BaseLayout) GetOnPostLayout() *LayoutEvent {
	return &this.OnPostLayout
}
