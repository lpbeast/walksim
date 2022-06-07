package itm

type Item struct {
	Name     string
	Desc     string
	Gettable bool
	Hasinv   bool
	Inv      map[string]Item
	Onuse    func(*Item) string
	Users    map[string]func(*Item) string
}

func (i Item) String() string {
	return i.Desc
}

func (i *Item) Use() string {
	return i.Onuse(i)
}

func (i *Item) Useon(target *Item) string {
	instr, ok := target.Users[i.Name]
	if ok {
		return instr(i)
	}
	return "You can't use " + i.Name + " on " + target.Name + ".\n"
}

//standard proc functions that items can call when used