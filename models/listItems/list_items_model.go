package listItems

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/zovenor/tea-models/models"
	"github.com/zovenor/tea-models/models/base"
)

type ListItemsModel struct {
	configs *Configs

	items     []*ListItemModel
	cursor    int
	findValue string
	findModel *textinput.Model
}

func NewListItemsModel(configs *Configs) *ListItemsModel {
	if configs == nil {
		configs = new(Configs)
	}
	_ = configs.check()
	lism := &ListItemsModel{
		configs: configs,
	}
	return lism
}

func (lism *ListItemsModel) AddItem(lim *ListItemModel) {
	lism.items = append(lism.items, lim)
	lism.setItemList()
}

func (lism *ListItemsModel) setItemList() {
	newItemList := make([]*ListItemModel, 0)
	groupsItems := make(map[string][]*ListItemModel, 0)
	deletedItems := make([]*ListItemModel, 0)
	index := 0
	for _, item := range lism.items {
		if item.deleted {
			deletedItems = append(deletedItems, item)
			continue
		}
		if item.group == "" {
			item.index = index
			newItemList = append(newItemList, item)
			index++
			continue
		}
		_, ok := groupsItems[item.group]
		if !ok {
			groupsItems[item.group] = make([]*ListItemModel, 0)
		}
		groupsItems[item.group] = append(groupsItems[item.group], item)
	}

	for _, groupItems := range groupsItems {
		for _, item := range groupItems {
			item.index = index
			newItemList = append(newItemList, item)
			index++
		}
	}

	for _, deletedItem := range deletedItems {
		deletedItem.index = index
		newItemList = append(newItemList, deletedItem)
		index++
	}
	lism.items = newItemList
}

func (lism *ListItemsModel) Configs() *Configs {
	return lism.configs
}

func (lism *ListItemsModel) Path() string {
	return fmt.Sprintf("%v > %v", lism.configs.ParentPath, lism.configs.Name)
}

func (lism *ListItemsModel) Items() []*ListItemModel {
	return lism.findItems(lism.findValue)
}

func (lism *ListItemsModel) findItems(findStr string) []*ListItemModel {
	if findStr == "" {
		return lism.items
	}
	needfulItems := make([]*ListItemModel, 0)
	for _, item := range lism.items {
		if strings.Contains(strings.ToLower(item.name), strings.ToLower(findStr)) {
			needfulItems = append(needfulItems, item)
		}
	}
	return needfulItems
}

func (lism *ListItemsModel) nextCursor() {
	for _, item := range lism.Items() {
		if lism.Cursor() < item.index {
			lism.cursor = item.index
			return
		}
	}
}

func (lism *ListItemsModel) lastCursor() {
	items := lism.Items()

	for i := len(items) - 1; i >= 0; i-- {
		if lism.cursor > items[i].index {
			lism.cursor = items[i].index
			return
		}
	}
}

func (lism *ListItemsModel) Page() int {
	cursor := lism.cursor
	maxPageItems := int(lism.configs.MaxPageItems)
	page := cursor / maxPageItems
	return page
}

func (lism *ListItemsModel) listItemsInPage() (pageItems []*ListItemModel, page int) {
	items := lism.Items()
	maxPageItems := int(lism.configs.MaxPageItems)
	if len(items) < maxPageItems {
		return items, 0
	}
	page = lism.Page()

	start := page * maxPageItems
	if start+maxPageItems > len(items) {
		pageItems = items[start:]
	} else {
		pageItems = items[start : start+maxPageItems]
	}

	return pageItems, page
}

func (lism *ListItemsModel) Cursor() int {
	items := lism.Items()
	exists := false
	for _, item := range items {
		if item.index == lism.cursor {
			exists = true
			break
		}
	}
	if !exists {
		newCursor := -1
		for _, item := range items {
			if item.index < lism.cursor {
				newCursor = item.index
			} else if item.index > lism.cursor {
				if newCursor == -1 {
					lism.cursor = item.index
				}
				break
			}
		}
		if newCursor != -1 {
			lism.cursor = newCursor
		} else {
			lism.cursor = 0
		}
	}
	if lism.cursor < 0 {
		lism.cursor = 0
	}
	return lism.cursor
}

func (lism *ListItemsModel) GetDeletedItems() []*ListItemModel {
	deletedItems := make([]*ListItemModel, 0)
	for _, item := range lism.items {
		if item.deleted {
			deletedItems = append(deletedItems, item)
		}
	}
	return deletedItems
}

func (lism *ListItemsModel) CurrentItem() *ListItemModel {
	if len(lism.items) == 0 {
		return nil
	}
	cursor := lism.Cursor()
	return lism.items[cursor]
}

func (lism *ListItemsModel) nextPage() {
	items := lism.Items()
	cursor := lism.Cursor()
	maxPageItems := int(lism.configs.MaxPageItems)
	page := lism.Page()
	if cursor+maxPageItems < len(items) {
		lism.cursor = cursor + maxPageItems
	} else if cursor < (page+1)*maxPageItems && !((page+1)*maxPageItems > len(items)) {
		lism.cursor = (page + 1) * maxPageItems
	}
}

func (lism *ListItemsModel) lastPage() {
	page := lism.Page()
	if page == 0 {
		return
	}
	cursor := lism.Cursor()
	maxPageItems := int(lism.configs.MaxPageItems)
	if cursor-maxPageItems < 0 {
		lism.cursor = 0
	} else {
		lism.cursor = cursor - maxPageItems
	}
}

func (lism *ListItemsModel) deleteSelectedItems() {
	if !lism.configs.DeletedMode {
		newItemList := make([]*ListItemModel, 0)
		for _, item := range lism.items {
			if !item.selected {
				newItemList = append(newItemList, item)
			}
		}
		lism.items = newItemList
	} else {
		for _, item := range lism.items {
			if item.selected {
				item.deleted = !item.deleted
				item.selected = false
			}
		}
	}
	lism.setItemList()
}

// Base methods for bubbletea models

func (lism *ListItemsModel) Init() tea.Cmd {
	return nil
}

func (lism *ListItemsModel) AllPages() int {
	items := lism.Items()
	if len(items)%int(lism.configs.MaxPageItems) != 0 {
		return len(items)/int(lism.configs.MaxPageItems) + 1
	} else {
		return len(items) / int(lism.configs.MaxPageItems)
	}
}

func (lism *ListItemsModel) View() string {
	var view string = fmt.Sprintf("%v\n\n", lism.Path())

	pageItems, page := lism.listItemsInPage()

	for _, lim := range pageItems {

		if lim.index == lism.Cursor() {
			view += fmt.Sprintf("%v ", lism.configs.CursorSymbol)
		} else {
			view += base.RepeatSymbol(" ", len(lism.configs.CursorSymbol)+1)
		}

		if lism.configs.SelectMode {
			if lim.selected {
				view += "[*] "
			} else {
				view += "[ ] "
			}
		}

		view += fmt.Sprintf("%v", lim.GetName())

		if lim.group != "" {
			view += color.New(color.FgHiCyan).Sprintf(" (%v)", lim.group)
		}

		if lim.deleted {
			view += color.New(color.FgRed).Sprint(" (deleted)")
		}

		view += "\n"
	}

	if lism.findModel != nil {
		view += fmt.Sprintf("\n\n%v\n\n", lism.findModel.View())
	}

	view += fmt.Sprintf("\n\nPage %v/%v\n\n", page+1, lism.AllPages())

	allKeys := make([]interface{}, 0)

	if lism.findValue != "" {
		if lism.findModel == nil {
			allKeys = append(allKeys, base.FindKey)
		} else {
			allKeys = append(allKeys, base.CancelKey, base.EnterKey)
		}
	}
	if lism.configs.SelectMode {
		allKeys = append(allKeys, base.SelectKey, base.DeleteKey)
	}

	allKeys = append(allKeys, base.ExitKey)
	for _, k := range lism.configs.Keys {
		allKeys = append(allKeys, k)
	}
	view += base.GetHints(allKeys...)

	return view
}

func (lism *ListItemsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if lism.configs.UpdateFunc != nil {
		model, cmd := (*lism.configs.UpdateFunc)(lism, msg)
		if model != nil || cmd != nil {
			return model, cmd
		}
	}

	if lism.findModel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case base.EnterKey:
				lism.findValue = lism.findModel.Value()
				lism.findModel = nil
				lism.setItemList()
			case base.CancelKey:
				lism.findModel = nil
			case base.ExitKey:
				return lism, tea.Quit
			default:
				var cmd tea.Cmd
				*lism.findModel, cmd = lism.findModel.Update(msg)
				return lism, cmd
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case base.ExitKey:
			return lism, tea.Quit
		case base.BackKey:
			if lism.configs.Parent != nil {
				return lism.configs.Parent, nil
			}
		case base.ForwardKey:
			if model, ok := base.IsForwardType(lism.CurrentItem().GetValue()); ok {
				return model, nil
			}
		case base.DownKey:
			lism.nextCursor()
		case base.UpKey:
			lism.lastCursor()
		case base.SelectKey:
			ci := lism.CurrentItem()
			ci.selected = !ci.selected
		case base.DeleteKey:
			confirmModel, err := models.NewConfirmModel("Do you want to delete selected items?", lism, lism.deleteSelectedItems)
			if err == nil {
				return confirmModel, nil
			}
		case base.FindKey:
			ti := textinput.New()
			ti.Placeholder = lism.findValue
			ti.SetValue(lism.findValue)
			ti.Focus()
			lism.findModel = &ti
			return lism, nil
		case "]":
			lism.nextPage()
		case "[":
			lism.lastPage()
		}
	}

	lism.Cursor()

	return lism, nil
}
