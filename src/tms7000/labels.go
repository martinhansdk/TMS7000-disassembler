package tms7000

type Label struct {
	address uint
	length  uint
	name    string
}

type LabelTable struct {
	labels map[uint]*Label
}

func NewLabelTable() *LabelTable {
	return &LabelTable{
		labels: make(map[uint]*Label),
	}
}

func (self *LabelTable) Add(address, length uint, name string) {
	self.labels[address] = &Label{
		address: address,
		length:  length,
		name:    name,
	}
}

func (self *LabelTable) Get(address uint) *Label {
	return self.labels[address]
}
