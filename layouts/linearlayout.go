package layouts

import (
	"log"

	"github.com/zzl/goforms/framework/utils"
	"github.com/zzl/goforms/layouts/aligns"
)

// LinearItem
// zero values as null values
type LinearItem struct {
	CollapsibleObject

	Control  Control
	ItemName string
	Name     string //control name

	Padding       int
	PaddingLeft   int
	PaddingTop    int
	PaddingRight  int
	PaddingBottom int

	Width    int
	MinWidth int

	Height    int
	MinHeight int

	Weight    float32
	Align     int
	SizeGroup string

	//
	Layout       Layout //sub layout
	ItemDefaults *LinearItem
	Items        []*LinearItem
}

func (this *LinearItem) GetControl() Control {
	return this.Control
}

func (this *LinearItem) GetName() string {
	return this.Name
}

func (this *LinearItem) GetLayout() Layout {
	return this.Layout
}

func (this *LinearItem) SetWidth(value int) {
	this.Width = value
}

func (this *LinearItem) SetHeight(value int) {
	this.Height = value
}

func (me LinearItem) GetItems() []LayoutItem {
	var items []LayoutItem
	for _, it := range me.Items {
		items = append(items, it)
	}
	return items
}

type LinearLayout struct {
	BaseLayout

	Items        []*LinearItem
	ItemDefaults *LinearItem
	Vertical     bool
	ContentAlign int //cross

	_items        []*LinearItem
	_sizeGroupMap map[string]int

	DebugName string

	analysisInfo *layoutAnalysisInfo
	bounds       Rect
}

type layoutLine struct {
	lineItems         []layoutLineItem
	sumAxisSize       int
	maxCrossSize      int
	assignedCrossSize int
	sumWeight         float32
	lastFlexIndex     int
}

type layoutLineItem struct {
	item             *LinearItem
	axisPadding      int
	axisStartPadding int
	axisEndPadding   int

	crossPadding      int
	crossStartPadding int
	crossEndPadding   int

	axisSize  int //preferred size
	crossSize int
	weight    float32
}

type layoutAnalysisInfo struct {
	layoutWidth  int
	layoutHeight int

	line layoutLine

	collapsedItems []*LinearItem
}

func (this *LinearLayout) Update() {
	this.analysisInfo = nil
	this.SetBoundsRect(this.GetBounds())
}

func (this *LinearLayout) GetBounds() Rect {
	return this.bounds
}

func (this *LinearLayout) Clone() Layout {
	clone := &LinearLayout{}
	if this.Items != nil {
		clone.Items = make([]*LinearItem, len(this.Items))
		copy(clone.Items, this.Items)
	}
	if this.ItemDefaults != nil {
		clone.ItemDefaults = &(*this.ItemDefaults)
	}
	return clone
}

func (this *LinearLayout) SetItemDefaults(itemDefaults LayoutItem) {
	this.ItemDefaults = itemDefaults.(*LinearItem)
}

func (this *LinearLayout) AddItems(items []LayoutItem, prepend bool) {
	var lItems []*LinearItem
	for _, item := range items {
		lItems = append(lItems, item.(*LinearItem))
	}
	if prepend {
		this.Items = append(lItems, this.Items...)
	} else {
		this.Items = append(this.Items, lItems...)
	}
}

func (this *LinearLayout) SetContainer(container Container) {
	items := this._getItems() //?
	for n, _ := range items {
		item := items[n]

		if this.Vertical {
			if item.Align == aligns.Default {
				item.Align = aligns.Stretch
			}
		}
		this.ResolveItemPaddings(item)

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
				item.Layout.AddItems(item.GetItems(), true)
			}
		}

		if item.Name != "" {
			if item.Control != nil {
				log.Fatal("??")
			}
			item.Control = container.GetControlByName(item.Name)
			if item.Control == nil {
				log.Fatal("layout item control not found: ", item.Name)
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
	this._sizeGroupMap = make(map[string]int)
	this.collectSizeGroupMap()
}

func (this *LinearLayout) SetSizeGroup(sg map[string]int) {
	this._sizeGroupMap = sg
	this.collectSizeGroupMap()
}

func (this *LinearLayout) ResolveItemPaddings(item *LinearItem) {
	//if item.Padding != 0 {
	//	println("?")
	//}
	utils.AssignDefault(&item.PaddingLeft, item.Padding)
	utils.AssignDefault(&item.PaddingTop, item.Padding)
	utils.AssignDefault(&item.PaddingRight, item.Padding)
	utils.AssignDefault(&item.PaddingBottom, item.Padding)
	utils.MagicZeroTo0(&item.PaddingLeft, &item.PaddingTop,
		&item.PaddingRight, &item.PaddingBottom)
}

//

func (this *LinearLayout) collectSizeGroupMap() {
	sgMap := this._sizeGroupMap
	for _, item := range this._items {
		if item.SizeGroup != "" {
			size, _ := sgMap[item.SizeGroup]
			cx := 0
			if item.Control != nil {
				cx, _ = item.Control.GetPreferredSize(4096, 4096)
			} else if item.Layout != nil { //?
				cx, _ = item.Layout.GetPreferredSize(4096, 4096)
			}
			if cx > size {
				sgMap[item.SizeGroup] = cx
			}
		}
		if item.Layout != nil {
			subLayout := item.Layout
			subLayout.SetSizeGroup(sgMap)
		}
	}
}

func applyLinearItemDefaults(item *LinearItem, itemDefaults *LinearItem) {
	i, d := item, itemDefaults
	utils.AssignDefault(&i.Align, d.Align)
	utils.AssignDefault(&i.Padding, d.Padding)
	utils.AssignDefault(&i.PaddingLeft, d.PaddingLeft)
	utils.AssignDefault(&i.PaddingTop, d.PaddingTop)
	utils.AssignDefault(&i.PaddingRight, d.PaddingRight)

	utils.AssignDefault(&i.PaddingBottom, d.PaddingBottom)
	utils.AssignDefault(&i.Width, d.Width)
	utils.AssignDefault(&i.MinWidth, d.MinWidth)
	utils.AssignDefault(&i.Height, d.Height)
	utils.AssignDefault(&i.MinHeight, d.MinHeight)
	utils.AssignDefault(&i.Weight, d.Weight)
	utils.AssignDefault(&i.SizeGroup, d.SizeGroup)

	//
	if i.Layout == nil && d.Layout != nil {
		i.Layout = d.Layout.Clone()
	}
	if i.ItemDefaults == nil && d.ItemDefaults != nil { //merge?
		i.ItemDefaults = d.ItemDefaults
	}
	if i.Items == nil {
		i.Items = make([]*LinearItem, len(d.Items))
		copy(i.Items, d.Items)
	}
}

func (this *LinearLayout) _getItems() []*LinearItem {
	items := this._items
	if items != nil {
		return items
	}
	if this.ItemDefaults != nil {
		itemCount := len(this.Items)
		items = make([]*LinearItem, itemCount)
		for n := 0; n < itemCount; n++ {
			items[n] = &*this.Items[n]
			applyLinearItemDefaults(this.Items[n], this.ItemDefaults)
		}
	} else {
		items = this.Items
	}
	this._items = items
	return items
}

func (this *LinearLayout) Analysis(layoutWidth int, layoutHeight int) *layoutAnalysisInfo {

	pInfo := this.analysisInfo
	if pInfo != nil && pInfo.layoutWidth == layoutWidth &&
		pInfo.layoutHeight == layoutHeight {
		return pInfo
	}

	vert := this.Vertical
	items := this._getItems()

	var layoutAxisSize, layoutCrossSize int
	if vert {
		layoutAxisSize, layoutCrossSize = layoutHeight, layoutWidth
	} else {
		layoutAxisSize, layoutCrossSize = layoutWidth, layoutHeight
	}

	var line layoutLine
	var collapsedItems []*LinearItem

	for n, it := range items {
		if layoutWidth == 1024 && layoutHeight == 0 || it.Collapsed {
			collapsedItems = append(collapsedItems, it)
			continue
		}
		_ = n
		var li layoutLineItem
		li.item = it
		//if li.item.ItemName == "XX" {
		//	println("?")
		//}
		if vert {
			li.axisStartPadding, li.axisEndPadding,
				li.crossStartPadding, li.crossEndPadding =
				it.PaddingTop, it.PaddingBottom, it.PaddingLeft, it.PaddingRight
		} else {
			li.axisStartPadding, li.axisEndPadding,
				li.crossStartPadding, li.crossEndPadding =
				it.PaddingLeft, it.PaddingRight, it.PaddingTop, it.PaddingBottom
		}
		li.axisPadding = li.axisStartPadding + li.axisEndPadding
		li.crossPadding = li.crossStartPadding + li.crossEndPadding

		control := it.Control
		cx, cy := 0, 0

		availableAxisSize := layoutAxisSize - (line.sumAxisSize + li.axisPadding)
		availableCrossSize := layoutCrossSize - li.crossPadding

		var availableWidth, availableHeight int
		if vert {
			availableWidth = availableCrossSize
			availableHeight = availableAxisSize
		} else {
			availableWidth = availableAxisSize
			availableHeight = availableCrossSize
		}

		if control != nil {
			cx, cy = control.GetPreferredSize(availableWidth, availableHeight)
		} else if it.Layout != nil {
			cx, cy = it.Layout.GetPreferredSize(availableWidth, availableHeight)
		}

		if it.Width != 0 {
			cx = it.Width
		} else if cx < it.MinWidth {
			cx = it.MinWidth
		}
		if it.Height != 0 {
			cy = it.Height
		} else if cy < it.MinHeight {
			cy = it.MinHeight
		}
		utils.MagicZeroTo0(&cx, &cy)

		if it.SizeGroup != "" { //vert?
			sgCx, ok := this._sizeGroupMap[it.SizeGroup]
			if ok {
				cx = sgCx
			}
		}
		if vert {
			li.axisSize = cy
			li.crossSize = cx
		} else {
			li.axisSize = cx
			li.crossSize = cy
		}

		//
		lineCrossSize := li.crossSize + li.crossPadding
		if lineCrossSize > line.maxCrossSize {
			line.maxCrossSize = lineCrossSize
		}

		if it.Weight == 0 {
			line.sumAxisSize += li.axisSize + li.axisPadding
		} else {
			li.weight = it.Weight
			line.sumWeight += it.Weight
			li.axisSize = -1
			line.lastFlexIndex = len(line.lineItems)
		}

		line.lineItems = append(line.lineItems, li)
	}

	//
	var info layoutAnalysisInfo
	info.line = line
	info.collapsedItems = collapsedItems
	return &info
}

//

func (this *LinearLayout) GetPreferredSize(layoutWidth int, layoutHeight int) (int, int) {

	info := this.Analysis(layoutWidth, layoutHeight)

	//
	axisSize := info.line.sumAxisSize
	crossSize := info.line.maxCrossSize

	if this.Vertical {
		return crossSize, axisSize
	} else {
		return axisSize, crossSize
	}
}

func (this *LinearLayout) SetBounds(left, top, width, height int) {
	this.SetBoundsRect(Rect{
		Left: left, Top: top, Right: left + width, Bottom: top + height})
}

func (this *LinearLayout) SetBoundsRect(bounds Rect) {

	info := this.Analysis(bounds.Width(), bounds.Height())

	//if this.DebugName == "TB" {
	//	println("?")
	//}

	//
	vert := this.Vertical
	var layoutAxisStart, layoutCrossStart int
	var layoutAxisSize, layoutCrossSize int
	if vert {
		layoutAxisStart, layoutCrossStart = bounds.Top, bounds.Left
		layoutAxisSize, layoutCrossSize = bounds.Height(), bounds.Width()
	} else {
		layoutAxisStart, layoutCrossStart = bounds.Left, bounds.Top
		layoutAxisSize, layoutCrossSize = bounds.Width(), bounds.Height()
	}

	line := &info.line
	line.assignedCrossSize = line.maxCrossSize

	var lineCrossStart int
	contentAlign := this.ContentAlign
	if contentAlign == aligns.Default {
		contentAlign = DefaultContentAlign
	}
	switch contentAlign {
	case aligns.Top:
		lineCrossStart = layoutCrossStart
	case aligns.Center:
		lineCrossStart = layoutCrossStart +
			(layoutCrossSize-line.maxCrossSize)/2
	case aligns.Bottom:
		lineCrossStart = layoutCrossStart + (layoutCrossSize - line.maxCrossSize)
	case aligns.Stretch:
		lineCrossStart = layoutCrossStart
		line.assignedCrossSize = layoutCrossSize
	}

	var controls []Control

	var flexSize, usedFlexSize int
	flexSize = layoutAxisSize - line.sumAxisSize
	axisStart := layoutAxisStart
	for n, li := range line.lineItems {
		//if li.item.ItemName == "XX" {
		//	println("?")
		//}
		axisStart += li.axisStartPadding
		axisSize := li.axisSize
		if axisSize == -1 {
			if n == line.lastFlexIndex {
				axisSize = flexSize - usedFlexSize
			} else {
				axisSize = utils.Round[float32, int](
					float32(flexSize) * li.weight / line.sumWeight)
			}
			usedFlexSize += axisSize
			axisSize -= li.axisPadding
		}

		crossStart := lineCrossStart
		crossSize := li.crossSize
		it := li.item
		align := it.Align
		if align == aligns.Default {
			align = DefaultItemAlign
		}
		switch align {
		case aligns.Top:
			crossStart += li.crossStartPadding
		case aligns.Center:
			crossStart += li.crossStartPadding +
				(line.assignedCrossSize-li.crossPadding-crossSize)/2
		case aligns.Bottom:
			crossStart += line.assignedCrossSize - (crossSize + li.crossEndPadding)
		case aligns.Stretch:
			crossStart += li.crossStartPadding
			crossSize = line.assignedCrossSize - li.crossPadding
		}

		var itemBounds Rect
		if vert {
			itemBounds = Rect{crossStart, axisStart,
				crossStart + crossSize, axisStart + axisSize}
		} else {
			itemBounds = Rect{axisStart, crossStart,
				axisStart + axisSize, crossStart + crossSize}
		}

		var ba BoundsAware
		if it.Control != nil {
			ba = it.Control
			controls = append(controls, it.Control)
		} else if it.Layout != nil {
			ba = it.Layout
		} else {
			ba = nil
		}
		if ba != nil {
			ba.SetBounds(itemBounds.Left, itemBounds.Top,
				itemBounds.Width(), itemBounds.Height())
		}

		axisStart += axisSize + li.axisEndPadding
	}
	lineCrossStart += line.assignedCrossSize

	//
	for _, c := range controls {
		c.Refresh()
	}

	//
	for _, it := range info.collapsedItems {
		if it.Control != nil {
			it.Control.SetBounds(0, 0, 1024, 0)
		} else if it.Layout != nil {
			it.Layout.SetBounds(0, 0, 1024, 0)
		}
	}

	//
	info.layoutWidth = bounds.Width()
	info.layoutHeight = bounds.Height()
	this.analysisInfo = info

	this.bounds = bounds
}

func (this *LinearLayout) FindItemByControl(control Control) LayoutItem {
	items := this._getItems()
	for n := 0; n < len(items); n++ {
		item := items[n]
		if item.Control == control {
			return item
		}
	}
	return nil
}

func (this *LinearLayout) GetItem(name string) LayoutItem {
	items := this._getItems()
	for n := 0; n < len(items); n++ {
		item := items[n]
		if item.ItemName == name {
			return item
		}
	}
	return nil
}
