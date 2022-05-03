package controller

type viewData struct {
	esiLink string
	data    interface{}
}

type iViewData interface {
	link() string
	setData(interface{})
}

func (x *viewData) setData(d interface{}) {
	x.data = d
}

func (x *viewData) link() string {
	return x.esiLink
}
func (x viewData) new() *viewData {
	return &x
}
