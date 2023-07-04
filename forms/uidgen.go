package forms

const (
	MenuGenIdStart    = -10000
	AccelGenIdStart   = -20000
	ToolbarGenIdStart = -30000
)

type UidGen struct {
	initId      int
	nextId      int
	step        int
	recycledIds []int
}

func NewUidGen(initId int, step int) *UidGen {
	return &UidGen{
		initId: initId,
		nextId: initId,
		step:   step,
	}
}

func (this *UidGen) Gen() int {
	var id int
	if len(this.recycledIds) > 0 {
		id = this.recycledIds[0]
		this.recycledIds = this.recycledIds[1:]
	} else {
		id = this.nextId
		this.nextId += this.step
	}
	return id
}

func (this *UidGen) Recycle(uid int) {
	this.recycledIds = append(this.recycledIds, uid)
}

func (this *UidGen) Reset() {
	this.nextId = this.initId
	this.recycledIds = nil
}

func (this *UidGen) IsGenerated(id int) bool {
	if this.step > 0 {
		return id >= this.initId && id < this.nextId
	} else {
		return id <= this.initId && id > this.nextId
	}
}
