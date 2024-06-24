package confirm

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/tea-models/models/base"
)

type ConfirmModel struct {
	description string
	f           func()
	parent      tea.Model
	actionKeys  base.ActionKeys
}

func (c *ConfirmModel) GetDescription() string { return c.description }

func NewConfirmModel(description string, parent tea.Model, f func()) (*ConfirmModel, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent is nil")
	}
	return &ConfirmModel{
		description: description,
		parent:      parent,
		f:           f,
	}, nil
}
func (c *ConfirmModel) SetActionKeys(keys base.ActionKeys) {
	c.actionKeys = keys
}

func (c *ConfirmModel) Init() tea.Cmd { return nil }

func (c *ConfirmModel) View() string {
	var s string

	s += c.description

	s += c.actionKeys.GetBaseHints(base.ExitKeyType, base.ConfirmKeyType, base.CancelKeyType)
	return s
}

func (c *ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch c.actionKeys.GetKeyTypeByHotKeyString(msg.String()) {
		case base.ExitKeyType:
			return c, tea.Quit
		case base.ConfirmKeyType:
			c.f()
			return c.parent, nil
		case base.CancelKeyType:
			return c.parent, nil
		}
	}
	return c, nil
}
