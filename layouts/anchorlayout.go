package layouts

import (
	"github.com/zzl/goforms/framework/consts"
	"github.com/zzl/goforms/framework/utils"
	"log"
	"math"
)

type AnchorItem struct {
	CollapsibleObject

	Control Control
	Name    string //control name

	Width    int
	MinWidth int

	Height    int
	MinHeight int

	//0=null,zero=0
	All    int
	Left   int
	Top    int
	Right  int
	Bottom int

	AnchorControl *AnchorControl

	Layout       Layout
	ItemDefaults LayoutItem
	Items        []LayoutItem
}

type AnchorControl struct {
	Left   Control
	Top    Control
	Right  Control
	Bottom Control
}

func (this *AnchorItem) GetControl() Control {
	return this.Control
}

func (this *AnchorItem) GetName() string {
	return this.Name
}

func (this *AnchorItem) GetLayout() Layout {
	return this.Layout
}

func (this *AnchorItem) SetWidth(value int) {
	this.Width = value
}

func (this *AnchorItem) SetHeight(value int) {
	this.Height = value
}

func (this *AnchorItem) GetItems() []LayoutItem {
	return this.Items
}

type AnchorLayout struct {
	BaseLayout

	Items        []*AnchorItem
	ItemDefaults *AnchorItem

	_items []*AnchorItem
	bounds Rect
}

type AnchorFlag byte

const (
	Data_AnchorFlag    = "__AnchorFlag"
	Data_AnchorControl = "__AnchorToControl"

	AnchorLeft   AnchorFlag = 0x1
	AnchorTop    AnchorFlag = 0x2
	AnchorRight  AnchorFlag = 0x4
	AnchorBottom AnchorFlag = 0x8
	AnchorAll    AnchorFlag = AnchorLeft | AnchorRight | AnchorTop | AnchorBottom
)

func (me AnchorFlag) String() string {
	s := ""
	if me&AnchorTop != 0 {
		s += ",Top"
	}
	if me&AnchorBottom != 0 {
		s += ",Bottom"
	}
	if me&AnchorLeft != 0 {
		s += ",Left"
	}
	if me&AnchorRight != 0 {
		s += ",Right"
	}
	if s == "" {
		s = "None"
	} else {
		s = s[1:]
	}
	return s
}

func SetAnchorFlag(ctrl Control, anchorFlag AnchorFlag) {
	//if anchorFlag & AnchorRight == 0 && anchorFlag & AnchorLeft == 0 {
	//	anchorFlag |= AnchorLeft
	//}
	//if anchorFlag & AnchorBottom == 0 && anchorFlag & AnchorTop == 0 {
	//	anchorFlag |= AnchorTop
	//}
	ctrl.SetData(Data_AnchorFlag, anchorFlag)
}

func SetAnchorControl(ctrl Control, anchorCtrl AnchorControl) {
	ctrl.SetData(Data_AnchorControl, &anchorCtrl)
}

func NewAnchorLayout() *AnchorLayout {
	obj := &AnchorLayout{}
	return obj
}

func (this *AnchorLayout) AddControl(control Control,
	anchorLeft, anchorTop, anchorRight, anchorBottom bool) {
	bounds := control.GetBounds()
	item := &AnchorItem{Control: control}
	if anchorLeft {
		item.Left = bounds.Left
	}
	if anchorTop {
		item.Top = bounds.Top
	}
	if anchorRight {
		item.Right = bounds.Right
	}
	if anchorBottom {
		item.Bottom = bounds.Bottom
	}
	this.Items = append(this.Items, item)
}

func applyAnchorItemDefaults(item *AnchorItem, itemDefaults *AnchorItem) {
	i, d := item, itemDefaults

	utils.AssignDefault(&i.Width, d.Width)
	utils.AssignDefault(&i.MinWidth, d.MinWidth)
	utils.AssignDefault(&i.Height, d.Height)
	utils.AssignDefault(&i.MinHeight, d.MinHeight)

	utils.AssignDefault(&i.Left, d.Left)
	utils.AssignDefault(&i.Top, d.Top)
	utils.AssignDefault(&i.Right, d.Right)
	utils.AssignDefault(&i.Bottom, d.Bottom)

	//
	if i.Layout == nil && d.Layout != nil {
		i.Layout = d.Layout.Clone()
	}
	if i.ItemDefaults == nil && d.ItemDefaults != nil { //merge?
		i.ItemDefaults = d.ItemDefaults
	}
	if i.Items == nil {
		//d._items?
		i.Items = make([]LayoutItem, len(d.Items))
		copy(i.Items, d.Items)
	}
}

func (this *AnchorLayout) SetSizeGroup(sg map[string]int) {
	//?
}

func (this *AnchorLayout) _getItems() []*AnchorItem {
	items := this._items
	if items != nil {
		return items
	}
	if this.ItemDefaults != nil {
		itemCount := len(this.Items)
		items = make([]*AnchorItem, itemCount)
		for n := 0; n < itemCount; n++ {
			applyAnchorItemDefaults(this.Items[n], this.ItemDefaults)
		}
	} else {
		items = this.Items
	}
	for _, it := range items {
		utils.AssignDefault(&it.Left, it.All)
		utils.AssignDefault(&it.Top, it.All)
		utils.AssignDefault(&it.Right, it.All)
		utils.AssignDefault(&it.Bottom, it.All)
		utils.MagicZeroTo0(&it.Left, &it.Top, &it.Right, &it.Bottom)
	}
	this._items = items
	return items
}

func (this *AnchorLayout) GetPreferredSize(maxWidth int, maxHeight int) (int, int) {
	items := this._getItems()
	minX, minY, maxX, maxY := math.MaxInt32, math.MaxInt32, 0, 0
	for n, _ := range items {
		item := items[n]

		x1, y1, x2, y2 := this.checkItemBounds(item, maxWidth, maxHeight)

		if x1 < minX {
			minX = x1
		}
		if x2 > maxX {
			maxX = x2
		}
		if y1 < minY {
			minY = y1
		}
		if y2 > maxY {
			maxY = y2
		}
	}
	return maxX - minX, maxY - minY
}

func (this *AnchorLayout) SetBoundsRect(bounds Rect) {
	ei := &LayoutEventInfo{
		Bounds: bounds,
	}
	this.OnPreLayout.Fire(this, ei)

	this.bounds = bounds
	maxWidth := bounds.Width()
	maxHeight := bounds.Height()

	xStart := bounds.Left
	yStart := bounds.Top

	items := this._getItems()
	for n, _ := range items {
		item := items[n]
		x1, y1, x2, y2 := this.checkItemBounds(item, maxWidth, maxHeight)

		var ba BoundsAware
		if item.Control != nil {
			ba = item.Control
		} else if item.Layout != nil {
			ba = item.Layout
		} else {
			ba = nil
		}
		if ba != nil {
			ba.SetBounds(xStart+x1, yStart+y1, x2-x1, y2-y1)
		}
	}
	//
	this.OnPostLayout.Fire(this, ei)
}

func (this *AnchorLayout) checkItemBounds(item *AnchorItem,
	maxWidth int, maxHeight int) (int, int, int, int) {

	//name := item.Control.GetName()
	//println("@", name)
	//if item.Control.GetName() == "BTN_REMOVE" {
	//	println("?")
	//}

	var cx, cy int
	if item.Control != nil {
		cx, cy = item.Control.GetPreferredSize(maxWidth, maxHeight)
	} else if item.Layout != nil {
		cx, cy = item.Layout.GetPreferredSize(maxWidth, maxHeight)
	}
	if item.Width != 0 {
		cx = item.Width
	} else if cx < item.MinWidth {
		cx = item.MinWidth
	}
	if item.Height != 0 {
		cy = item.Height
	} else if cy < item.MinHeight {
		cy = item.MinHeight
	}

	var x1, y1, x2, y2 int
	l, t, r, b := item.Left, item.Top, item.Right, item.Bottom
	//if item.Control.GetName() == "XX" {
	//	println("?")
	//}
	if l != 0 && r != 0 {
		if l == consts.Zero {
			l = 0
		}
		if r == consts.Zero {
			r = 0
		}
		x1 = l
		if item.AnchorControl != nil && item.AnchorControl.Left != nil {
			x1 = item.AnchorControl.Left.GetBounds().Right + l
		}
		x2 = maxWidth - r
		if item.AnchorControl != nil && item.AnchorControl.Right != nil {
			x2 = item.AnchorControl.Right.GetBounds().Left - r
		}
	} else if r != 0 {
		if r == consts.Zero {
			r = 0
		}
		x2 = maxWidth - r
		if item.AnchorControl != nil && item.AnchorControl.Right != nil {
			x2 = item.AnchorControl.Right.GetBounds().Left - r
		}
		x1 = x2 - cx
	} else {
		if l == consts.Zero {
			l = 0
		}
		x1 = l
		if item.AnchorControl != nil && item.AnchorControl.Left != nil {
			x1 = item.AnchorControl.Left.GetBounds().Right + l
		}
		x2 = x1 + cx
	}
	if t != 0 && b != 0 {
		if t == consts.Zero {
			t = 0
		}
		if b == consts.Zero {
			b = 0
		}
		y1 = t
		if item.AnchorControl != nil && item.AnchorControl.Top != nil {
			y1 = item.AnchorControl.Left.GetBounds().Bottom + t
		}
		y2 = maxHeight - b
		if item.AnchorControl != nil && item.AnchorControl.Bottom != nil {
			y2 = item.AnchorControl.Bottom.GetBounds().Top - b
		}
	} else if b != 0 {
		if b == consts.Zero {
			b = 0
		}
		y2 = maxHeight - b
		if item.AnchorControl != nil && item.AnchorControl.Bottom != nil {
			y2 = item.AnchorControl.Bottom.GetBounds().Top - b
		}
		y1 = y2 - cy
	} else {
		if t == consts.Zero {
			t = 0
		}
		y1 = t
		if item.AnchorControl != nil && item.AnchorControl.Top != nil {
			y1 = item.AnchorControl.Left.GetBounds().Bottom + t
		}
		y2 = y1 + cy
	}
	return x1, y1, x2, y2
}

func (this *AnchorLayout) SetBounds(left, top, width, height int) {
	this.SetBoundsRect(Rect{left, top, left + width, top + height})
}

func (this *AnchorLayout) GetBounds() Rect {
	return this.bounds
}

func (this *AnchorLayout) SetItemDefaults(itemDefaults LayoutItem) {
	this.ItemDefaults = itemDefaults.(*AnchorItem)
}

func (this *AnchorLayout) SetContainer(container Container) {
	items := this._getItems()
	for n, _ := range items {
		item := items[n]
		//item.Layout
		if item.Items != nil || item.ItemDefaults != nil {
			if item.Layout == nil {
				item.Layout = &LinearLayout{
					DebugName: "(auto generated)",
				}
			}
			if item.ItemDefaults != nil {
				item.Layout.SetItemDefaults(item.ItemDefaults)
			}
			if item.Items != nil {
				item.Layout.AddItems(item.Items, true)
			}
		}

		if item.Name != "" {
			if item.Control != nil {
				log.Fatal("??")
			}
			item.Control = container.GetControlByName(item.Name)
			if item.Control == nil {
				log.Fatal("??")
			}
			la, ok := item.Control.(LayoutAware)
			if ok {
				la.SetLayout(item.Layout)
			}
		} else if item.Layout != nil {
			item.Layout.SetContainer(container)
		}

		if item.Control != nil {
			item.Control.SetData(Data_Layout, this)
		}
	}
}

func (this *AnchorLayout) AddItems(items []LayoutItem, prepend bool) {
	var aItems []*AnchorItem
	for _, item := range items {
		aItems = append(aItems, item.(*AnchorItem))
	}
	if prepend {
		//var thisItems = this.Items
		//if thisItems == nil {
		//	thisItems = make([]*AnchorItem, 0)
		//}
		this.Items = append(aItems, this.Items...)
	} else {
		this.Items = append(this.Items, aItems...)
	}
}

func (this *AnchorLayout) Clone() Layout {
	clone := &AnchorLayout{}
	if this.Items != nil {
		clone.Items = make([]*AnchorItem, len(this.Items))
		copy(clone.Items, this.Items)
	}
	if this.ItemDefaults != nil {
		clone.ItemDefaults = &(*this.ItemDefaults)
	}
	//return clone
	return nil
}

func (this *AnchorLayout) FindItemByControl(control Control) LayoutItem {
	items := this._getItems()
	for n := 0; n < len(items); n++ {
		item := items[n]
		if item.Control == control {
			return item
		}
	}
	return nil
}

func (this *AnchorLayout) Update() {
	this.SetBoundsRect(this.GetBounds())
}

func (this *AnchorLayout) ParseItems(container Container) {
	childWins := container.GetControls()
	parentWidth, parentHeight := container.GetClientSize()
	for _, childWin := range childWins {
		control, ok := childWin.(Control)
		if !ok {
			continue
		}
		//println("#", control.GetName())

		//if control.GetName() == "ALL_LB" {
		//	println("?")
		//}

		data := control.GetData(Data_AnchorFlag)
		if data == nil {
			continue
		}
		flag := data.(AnchorFlag)
		if flag == 0 {
			continue
		}

		if flag&AnchorRight == 0 && flag&AnchorLeft == 0 {
			flag |= AnchorLeft
		}
		if flag&AnchorBottom == 0 && flag&AnchorTop == 0 {
			flag |= AnchorTop
		}

		//
		var anchorControl *AnchorControl
		data = control.GetData(Data_AnchorControl)
		if data != nil {
			anchorControl = data.(*AnchorControl)
		}

		//println("..")

		//?
		//(*WindowObject)(control.GetObjectPointer()).unsetDluFlag()

		bounds := control.GetBounds()
		item := &AnchorItem{
			Control:       control,
			Left:          bounds.Left,
			Top:           bounds.Top,
			Width:         bounds.Width(),
			Height:        bounds.Height(),
			Right:         parentWidth - bounds.Right,
			Bottom:        parentHeight - bounds.Bottom,
			AnchorControl: anchorControl,
		}
		if flag&AnchorLeft == 0 {
			item.Left = 0
		} else {
			if anchorControl != nil && anchorControl.Left != nil {
				item.Left = bounds.Left - anchorControl.Left.GetBounds().Right
			}
			if item.Left == 0 {
				item.Left = consts.Zero
			}
		}
		if flag&AnchorTop == 0 {
			item.Top = 0
		} else {
			if anchorControl != nil && anchorControl.Top != nil {
				item.Top = bounds.Top - anchorControl.Top.GetBounds().Bottom
			}
			if item.Top == 0 {
				item.Top = consts.Zero
			}
		}
		if flag&AnchorRight == 0 {
			item.Right = 0
		} else {
			if anchorControl != nil && anchorControl.Right != nil {
				item.Right = anchorControl.Right.GetBounds().Left - bounds.Right
			}
			if item.Right == 0 {
				item.Right = consts.Zero
			}
		}
		if flag&AnchorBottom == 0 {
			item.Bottom = 0
		} else {
			if anchorControl != nil && anchorControl.Bottom != nil {
				item.Bottom = anchorControl.Bottom.GetBounds().Top - bounds.Bottom
			}
			if item.Bottom == 0 {
				item.Bottom = consts.Zero
			}
		}
		if item.Left != 0 && item.Right != 0 {
			item.Width = 0
		}
		if item.Top != 0 && item.Bottom != 0 {
			item.Height = 0
		}

		//
		this.Items = append(this.Items, item)
	}
	println("?")
}

func (this *AnchorLayout) GetItem(name string) LayoutItem {
	//todo..
	return nil
}
