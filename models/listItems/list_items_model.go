package listItems

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/zovenor/tea-models/models/base"
	"github.com/zovenor/tea-models/models/confirm"
)

type ListItemsModel struct {
	configs *Configs

	items        []*ListItemModel
	cursor       int
	findValue    string
	findModel    *textinput.Model
	windowParams base.WindowParams
}

func NewListItemsModel(configs *Configs, opts ...func(*Configs)) (*ListItemsModel, error) {
	if configs == nil {
		configs = new(Configs)
	}
	// Base opts
	opts = append([]func(*Configs){WithBaseConfigsView}, opts...)
	// Set opts
	for _, opt := range opts {
		opt(configs)
	}
	err := configs.check()
	if err != nil {
		return nil, err
	}
	lism := &ListItemsModel{
		configs: configs,
	}
	return lism, nil
}

func (lism *ListItemsModel) SetMapValue(key string, value string) {
	lism.configs.MapArgs[key] = value
}

func (lism *ListItemsModel) GetMapValue(key string) interface{} {
	return lism.configs.MapArgs[key]
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
					newCursor = item.index
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
	allPath := make([]string, 0, len(lism.configs.ParentPath)+1)
	allPath = append(allPath, lism.configs.ParentPath...)
	allPath = append(allPath, lism.configs.Name)

	view := lism.configs.ConfigsViewTheme.Title(allPath, lism.windowParams)
	pageItems, page := lism.listItemsInPage()

	for _, lim := range pageItems {
		active := lim.index == lism.Cursor()
		view += lism.Configs().ConfigsViewTheme.ItemView(lim, active, lism.windowParams)
	}
	var findValue *string
	if lism.findModel != nil {
		fv := lism.findModel.Value()
		findValue = &fv
	}
	view += lism.configs.ConfigsViewTheme.Footer(
		uint64(page),
		uint64(lism.AllPages()),
		lism.groupItems(),
		findValue,
		lism.windowParams,
	)

	return view
}

func (lism *ListItemsModel) groupItemsList() [][]*ListItemModel {
	groupsItems := make([][]*ListItemModel, 0)
GroupItemsLoop:
	for _, item := range lism.Items() {
		for i, groupItems := range groupsItems {
			if groupItems[0].group == item.group {
				groupsItems[i] = append(groupsItems[i], item)
				continue GroupItemsLoop
			}
		}
		newLI := make([]*ListItemModel, 1)
		newLI[0] = item
		groupsItems = append(groupsItems, newLI)
	}
	return groupsItems
}

func (lism *ListItemsModel) groupItems() []ItemsGroup {
	items := make([]ItemsGroup, 0)
	allItemsString := "all items"
	if allItemsNewValue, exists := lism.configs.RenameGroupsView["$allItems"]; exists {
		allItemsString = allItemsNewValue
	}

	items = append(items, ItemsGroup{
		Name:  allItemsString,
		Total: uint64(len(lism.Items())),
	})
	if lism.configs.MoreItemsLenInfo {
		for _, gItems := range lism.groupItemsList() {
			var groupNameView string
			if gItems[0].group != "" {
				groupNameView = fmt.Sprintf("%v", gItems[0].group)
			} else {
				if newKey, exists := lism.configs.RenameGroupsView[gItems[0].group]; exists {
					groupNameView = fmt.Sprintf("%v", newKey)
				}
			}
			if gItems[0].group == "" {
				if !lism.configs.GroupsView {
					continue
				}
				if _, exists := lism.configs.RenameGroupsView[gItems[0].group]; !exists {
					groupNameView = "no group"
				}
			}
			items = append(items, ItemsGroup{
				Name:  groupNameView,
				Total: uint64(len(gItems)),
			})
		}
	}
	return items
}

func (lism *ListItemsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Set window size
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		lism.windowParams = base.WindowParams{
			Width:  uint64(msg.Width),
			Height: uint64(msg.Height),
		}
	}
	if lism.findModel != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch lism.configs.ActionKeys.GetKeyTypeByHotKeyString(msg.String()) {
			case base.EnterKeyType:
				lism.findValue = lism.findModel.Value()
				lism.findModel = nil
				lism.setItemList()
			case base.CancelKeyType:
				lism.findModel = nil
			case base.ExitKeyType:
				return lism, tea.Quit
			default:
				var cmd tea.Cmd
				*lism.findModel, cmd = lism.findModel.Update(msg)
				return lism, cmd
			}
		}
	}

	if lism.configs.UpdateFunc != nil {
		model, cmd := (*lism.configs.UpdateFunc)(lism, msg)
		if model != nil || cmd != nil {
			return model, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch lism.configs.ActionKeys.GetKeyTypeByHotKeyString(msg.String()) {
		case base.ExitKeyType:
			return lism, tea.Quit
		case base.BackKeyType:
			if lism.configs.Parent != nil {
				return lism.configs.Parent, nil
			}
		case base.ForwardKeyType:
			if model, ok := base.IsForwardType(lism.CurrentItem().GetValue()); ok {
				return model, nil
			}
		case base.DownKeyType:
			lism.nextCursor()
		case base.UpKeyType:
			lism.lastCursor()
		case base.SelectKeyType:
			ci := lism.CurrentItem()
			ci.selected = !ci.selected
		case base.DeleteKeyType:
			confirmModel, err := confirm.NewConfirmModel("Do you want to delete selected items?", lism, lism.deleteSelectedItems)
			if err == nil {
				return confirmModel, nil
			}
		case base.FindKeyType:
			ti := textinput.New()
			ti.Placeholder = lism.findValue
			ti.SetValue(lism.findValue)
			ti.Focus()
			lism.findModel = &ti
			return lism, nil
		case base.NextPageKeyType:
			lism.nextPage()
		case base.PreviousPageKeyType:
			lism.lastPage()
		}
	}

	lism.cursor = lism.Cursor()

	return lism, nil
}
