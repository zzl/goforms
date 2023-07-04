package layouts

import (
	"github.com/zzl/goforms/framework/events"
	"github.com/zzl/goforms/framework/types"
)

type SimpleEvent = events.SimpleEvent
type SimpleEventInfo = events.SimpleEventInfo
type Rect = types.Rect
type BoundsAware = types.BoundsAware

type Container interface {
	GetControlByName(name string) Control
	GetControls() []Control
	GetClientSize() (cx, cy int)
}

type Control interface {
	types.NameAware
	types.DataAware
	BoundsAware

	GetPreferredSize(int, int) (cx, cy int)
	Refresh()
}
