package listItems

import "reflect"

type ListItemModel struct {
	name  string
	value reflect.Value
	group string

	selected bool
	index    int
	deleted  bool
}

func NewListItemModel() *ListItemModel {
	return new(ListItemModel)
}

func (lim *ListItemModel) GetGroup() string {
	return lim.group
}

func (lim *ListItemModel) SetGroup(group string) {
	lim.group = group
}

func (lim *ListItemModel) GetName() string {
	return lim.name
}

func (lim *ListItemModel) SetName(name string) {
	lim.name = name
}

func (lim *ListItemModel) GetValue() reflect.Value {
	return lim.value
}

func (lim *ListItemModel) SetValue(value interface{}) {
	lim.value = reflect.ValueOf(value)
}

func (lim *ListItemModel) SetDeleteStatus(deleted bool) {
	lim.deleted = deleted
}

func (lim *ListItemModel) GetDeleteStatus() bool {
	return lim.deleted
}
