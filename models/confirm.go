package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/tea-models/models/base"
)

type ConfirmModel struct {
	description string
	f           func()
	parent      tea.Model
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

func (c *ConfirmModel) Init() tea.Cmd { return nil }

func (c *ConfirmModel) View() string {
	var s string

	s += c.description

	s += base.GetHints(base.ExitKey, base.ConfirmKey, base.CancelKey)

	return s
}

func (c *ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case base.ExitKey:
			return c, tea.Quit
		case base.ConfirmKey:
			c.f()
			return c.parent, nil
		case base.CancelKey:
			return c.parent, nil
		}
	}
	return c, nil
}
