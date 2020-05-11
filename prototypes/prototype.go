package prototype

type Prototype struct {
	Style Style
	*Brief
	Menu *Menu
	*Content
	ItemList  *ItemList
	UserInput *UserInput
	Footer    *Item
}

func New() *Prototype {
	return &Prototype{}
}

type Brief struct {
	Icon        Icon
	Label       string
	Description string
	Target      string
	Tags        []Tag
}

func NewBreif() *Brief {
	return &Brief{}
}

type Content struct {
	Timestamp int64
	Cover     []string
	Body      string
	Related   []*Brief
}

func NewContent() *Content {
	return &Content{}
}

type Item struct {
	*Brief
	*Content
}

func NewItem() *Item {
	return &Item{
		Brief:   NewBreif(),
		Content: NewContent(),
	}
}

type ItemList struct {
	Items  []*Item
	Source API
}

func NewItemList() *ItemList {
	return &ItemList{}
}

type Menu struct {
	Menus []*Brief
}

func NewMenu() *Menu {
	return &Menu{}
}
