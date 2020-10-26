package storage

type Group struct {
	group []Storager
}

func NewGroup(storagers ...Storager) *Group {
	group := make([]Storager, 0, len(storagers))
	group = append(group, storagers...)
	return &Group{
		group: group,
	}
}

func (g *Group) Add(s Storager) {
	if s != nil {
		g.group = append(g.group, s)
	}
}

func (g *Group) Size() int {
	return len(g.group)
}

func (g *Group) Write(msg *Message) (err error) {
	for _, storager := range g.group {
		if ierr := storager.Write(msg); ierr != nil {
			err = ierr
		}
	}
	return
}

func (g *Group) Close() (err error) {
	for _, storager := range g.group {
		if ierr := storager.Close(); ierr != nil {
			err = ierr
		}
	}
	return
}
