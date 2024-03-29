package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/zovenor/tea-models/models/base"
)

type ItemMsg *ListItemModel
type ItemsMsg []*ListItemModel

type ListItemModel struct {
	name  string
	value interface{}
	group string

	selected bool
}

func (im *ListItemModel) SetGroup(group string) {
	im.group = group
}

func (im *ListItemModel) GetGroup() string {
	return im.group
}

func (im *ListItemModel) GetName() string {
	return im.name
}

func (im *ListItemModel) GetValue() interface{} {
	return im.value
}
func NewListItemModel(name string, value interface{}) *ListItemModel {
	return &ListItemModel{
		name:  name,
		value: value,
	}
}

type ListItemsConf struct {
	Name           string
	SelectMode     bool
	ReturnValue    bool
	FindMode       bool
	ParentPath     string
	Parent         tea.Model
	MaxItemsInPage int
	Indexes        bool
	KeyValues      map[string]interface{}
	CmdsF          []func(lim *ListItemsModel) tea.Cmd
	UpdateF        *func(*ListItemsModel, tea.Msg) (tea.Model, tea.Cmd)
	ErrForward     bool
}

func NewListItemsModel(listItemsConf ListItemsConf) (*ListItemsModel, error) {
	if listItemsConf.MaxItemsInPage < 1 {
		return nil, fmt.Errorf("maxItemsInPage should be more than 0")
	}
	lim := &ListItemsModel{
		name:        listItemsConf.Name,
		selectMode:  listItemsConf.SelectMode,
		returnValue: listItemsConf.ReturnValue,
		parentPath:  listItemsConf.ParentPath,

		parent: listItemsConf.Parent,

		cursorSymbol:   base.CursorSymbol,
		maxItemsInPage: listItemsConf.MaxItemsInPage,
		indexes:        listItemsConf.Indexes,
		keyValues:      listItemsConf.KeyValues,
		cmdsF:          listItemsConf.CmdsF,
		updateF:        listItemsConf.UpdateF,
		findMode:       listItemsConf.FindMode,
		errForward:     listItemsConf.ErrForward,
	}
	return lim, nil
}

type ListItemsModel struct {
	name        string
	items       []*ListItemModel
	selectMode  bool
	returnValue bool
	parentPath  string
	indexes     bool
	keyValues   map[string]interface{}
	cmdsF       []func(lim *ListItemsModel) tea.Cmd
	updateF     *func(*ListItemsModel, tea.Msg) (tea.Model, tea.Cmd)
	view        *string
	findMode    bool

	cursor               int
	cursorSymbol         string
	err                  error
	status               string
	parent               tea.Model
	viewListItemsIndexed []int
	findValue            string
	findModel            *textinput.Model
	findCursor           int
	maxItemsInPage       int
	keys                 []base.Key
	errForward           bool
}

func (lim *ListItemsModel) GetKeyValues() map[string]interface{} {
	return lim.keyValues
}

func (lim *ListItemsModel) GetValueByKey(key string) interface{} {
	for k, v := range lim.keyValues {
		if k == key {
			return v
		}
	}
	return nil
}

func (lim *ListItemsModel) SetKeyValueByKey(key string, value interface{}) {
	lim.keyValues[key] = value
}
func (lim *ListItemsModel) AddItem(name string, value interface{}) *ListItemModel {
	im := NewListItemModel(name, value)
	lim.items = append(lim.items, im)
	lim.filterByName(lim.findValue)
	lim.setCursorByFindCursor()
	return im
}

func (lim *ListItemsModel) GetItems() []*ListItemModel {
	return lim.items
}

func (lim *ListItemsModel) SetItems(newItems []*ListItemModel) {
	lim.items = newItems
	lim.filterByName(lim.findValue)
	lim.setCursorByFindCursor()
}

func (lim *ListItemsModel) SetCursorSymbol(cursorSymbol string) {
	lim.cursorSymbol = cursorSymbol
}
func (lim *ListItemsModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, cmdF := range lim.cmdsF {
		cmds = append(cmds, cmdF(lim))
	}
	return tea.Batch(cmds...)
}

func (lim *ListItemsModel) SetStatus(status string) {
	lim.status = status
}

type ListItemModelWithIndex struct {
	index int
	lim   *ListItemModel
}

func sortByGroup(items []*ListItemModelWithIndex) [][]*ListItemModelWithIndex {
	listItemGroups := make([][]*ListItemModelWithIndex, 0)
LimLoop:
	for _, limWithIndex := range items {
		for i := range listItemGroups {
			if listItemGroups[i][0].lim.group == limWithIndex.lim.group {
				listItemGroups[i] = append(listItemGroups[i], limWithIndex)
				continue LimLoop
			}
		}
		listItemGroups = append(listItemGroups, []*ListItemModelWithIndex{limWithIndex})
	}
	return listItemGroups
}

func (lim *ListItemsModel) View() string {
	if lim.view != nil {
		return *lim.view
	}
	var s string

	if lim.parentPath != "" {
		s += fmt.Sprintf("%v > ", lim.parentPath)
	}

	s += fmt.Sprintf("%v\n\n", lim.name)

	itemIndexes := lim.getPageItemsIndexes()
	items := make([]*ListItemModelWithIndex, 0)
	for _, index := range itemIndexes {
		item, err := lim.getItemByIndex(index)
		if err != nil {
			continue
		}
		items = append(items,
			&ListItemModelWithIndex{
				index: index,
				lim:   item,
			})
	}
	groupsItems := sortByGroup(items)

	for _, groupItems := range groupsItems {
		for i, item := range groupItems {

			if lim.cursor == item.index {
				s += lim.cursorSymbol + " "
			} else {
				s += base.RepeatSymbol(" ", len(lim.cursorSymbol)+1)
			}

			if lim.selectMode {
				if item.lim.selected {
					s += "[*] "
				} else {
					s += "[ ] "
				}
			}
			if lim.indexes {
				s += fmt.Sprintf("%v) ", i+1)
			}
			s += fmt.Sprintf("%v", item.lim.name)
			if item.lim.group != "" {
				s += fmt.Sprintf(" (%v)", item.lim.group)
			}
			s += "\n"
		}
		s += "\n"
	}

	if lim.findModel != nil {
		s += fmt.Sprintf("\n%v\n", lim.findModel.View())
	}

	if lim.err != nil {
		s += "\n"
		s += color.New(color.BgRed).Sprintf("error: %v", lim.err.Error())
		s += "\n"
		lim.err = nil
	}

	if lim.status != "" {
		s += "\n"
		s += color.New(color.BgBlue).Sprint(lim.status)
		s += "\n"
		lim.status = ""
	}

	allKeys := make([]interface{}, 0)

	if lim.findMode {
		if lim.findModel == nil {
			allKeys = append(allKeys, base.FindKey)
		} else {
			allKeys = append(allKeys, base.CancelKey, base.EnterKey)
		}
	}
	if lim.selectMode {
		allKeys = append(allKeys, base.SelectKey)
	}

	allKeys = append(allKeys, base.ExitKey)
	for _, k := range lim.keys {
		allKeys = append(allKeys, k)
	}
	s += base.GetHints(allKeys...)
	if lim.getPagesLen() > 0 {
		s += fmt.Sprintf("Page %v/%v. All items: %v", lim.getPageIndex()+1, lim.getPagesLen(), len(lim.viewListItemsIndexed))
	}

	return s
}

func (lim *ListItemsModel) getItemByIndex(index int) (*ListItemModel, error) {
	if index >= len(lim.items) {
		return nil, fmt.Errorf("index out of range")
	} else {
		return lim.items[index], nil
	}
}

func (lim *ListItemsModel) GetParent() (tea.Model, error) {
	if lim.parent == nil {
		return nil, fmt.Errorf("parent is nil")
	}
	return lim.parent, nil
}

func (lim *ListItemsModel) filterByName(s string) {
	var neededIndexes []int
	for i, item := range lim.items {
		if strings.Contains(strings.ToLower(item.name), strings.ToLower(s)) {
			neededIndexes = append(neededIndexes, i)
		}
	}
	lim.viewListItemsIndexed = neededIndexes
}

func (lim *ListItemsModel) setCursorByFindCursor() {
	if len(lim.viewListItemsIndexed) == 0 {
		lim.findCursor = 0
		lim.cursor = 0
		return
	}
	if lim.findCursor >= len(lim.viewListItemsIndexed) {
		lim.findCursor = len(lim.viewListItemsIndexed) - 1
	}
	if lim.viewListItemsIndexed[lim.findCursor] < len(lim.items) {
		lim.cursor = lim.viewListItemsIndexed[lim.findCursor]
	} else {
		lim.setCursorByFindCursor()
	}
}

func (lim *ListItemsModel) nextIndex() {
	nextFindCursor := lim.findCursor + 1
	if nextFindCursor < len(lim.viewListItemsIndexed) {
		lim.findCursor = nextFindCursor
		lim.setCursorByFindCursor()
	} else {
		lim.setCursorByFindCursor()
	}
}

func (lim *ListItemsModel) getPagesLen() int {
	if len(lim.viewListItemsIndexed) < lim.maxItemsInPage {
		return 1
	}
	return len(lim.viewListItemsIndexed)/lim.maxItemsInPage + 1
}

func (lim *ListItemsModel) getPageIndex() int {
	if len(lim.viewListItemsIndexed) == 0 {
		return 0
	}
	return lim.findCursor / lim.maxItemsInPage
}

func (lim *ListItemsModel) getPageItemsIndexes() []int {
	if len(lim.viewListItemsIndexed) < lim.maxItemsInPage {
		return lim.viewListItemsIndexed
	} else if lim.getPageIndex()+1 < lim.getPagesLen() {
		return lim.viewListItemsIndexed[lim.getPageIndex()*lim.maxItemsInPage : (lim.getPageIndex()+1)*lim.maxItemsInPage]
	} else {
		return lim.viewListItemsIndexed[lim.getPageIndex()*lim.maxItemsInPage:]
	}
}

func (lim *ListItemsModel) lastIndex() {
	lastFindCursor := lim.findCursor - 1
	if lastFindCursor >= 0 {
		lim.findCursor = lastFindCursor
	}
	lim.setCursorByFindCursor()
}

func (lim *ListItemsModel) getSelectedItemsMsg() ItemsMsg {
	var items ItemsMsg
	for _, item := range lim.items {
		if item.selected {
			items = append(items, item)
		}
	}
	return items
}

func (lim *ListItemsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if lim.updateF != nil {
		model, cmd := (*lim.updateF)(lim, msg)
		if model != nil || cmd != nil {
			return model, cmd
		}
	}

	if lim.findModel == nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case base.UpKey:
				lim.lastIndex()
			case base.DownKey:
				lim.nextIndex()
			case base.ForwardKey:
				neededItem, err := lim.getItemByIndex(lim.cursor)
				if err != nil {
					lim.err = err
					return lim, nil
				}
				switch value := neededItem.value.(type) {
				case tea.Model:
					return value, nil
				default:
					if lim.errForward {
						lim.err = fmt.Errorf("can not forward to this item")
					}
					return lim, nil
				}
			case base.BackKey:
				parent, err := lim.GetParent()
				if err != nil {
					if lim.errForward {
						lim.SetError(err)
					}
					return lim, nil
				}
				return parent.Update(nil)
			case base.EnterKey:
				if lim.returnValue {
					parent, err := lim.GetParent()
					if err != nil {
						lim.err = err
						return lim, nil
					}
					if !lim.selectMode {
						neededItem, err := lim.getItemByIndex(lim.cursor)
						if err != nil {
							lim.err = err
							return lim, nil
						}
						return parent.Update(ItemMsg(neededItem))
					} else {
						neededItems := lim.getSelectedItemsMsg()
						return parent.Update(neededItems)
					}
				}
			case base.SelectKey:
				if lim.selectMode {
					neededItem, err := lim.getItemByIndex(lim.cursor)
					if err != nil {
						lim.err = err
						return lim, nil
					}
					neededItem.selected = !neededItem.selected
				}
			case base.FindKey:
				if lim.findMode {
					ti := textinput.New()
					ti.Placeholder = lim.findValue
					ti.SetValue(lim.findValue)
					ti.Focus()
					lim.findModel = &ti
				}
			case base.ExitKey:
				return lim, tea.Quit
			}
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case base.EnterKey:
				lim.findValue = lim.findModel.Value()
				lim.filterByName(lim.findValue)
				lim.findModel = nil
				lim.setCursorByFindCursor()
			case base.CancelKey:
				lim.findModel = nil
			case base.ExitKey:
				return lim, tea.Quit
			default:
				var cmd tea.Cmd
				*lim.findModel, cmd = lim.findModel.Update(msg)
				return lim, cmd
			}
		}
	}

	return lim, nil
}

func (lim *ListItemsModel) SetError(err error) {
	lim.err = err
}

func (lim *ListItemsModel) SetView(view *string) {
	lim.view = view
}

func (lim *ListItemsModel) SetNewKeyForView(key string, description string) {
	lim.keys = append(lim.keys, base.Key{
		Name:        key,
		Description: description,
	})
}

func (lim *ListItemsModel) GetCurrentItem() (*ListItemModel, error) {
	return lim.getItemByIndex(lim.cursor)
}
