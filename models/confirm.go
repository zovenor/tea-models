package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/tea-models/models/base"
)

type ConfirmMsg struct {
	Value   interface{}
	Confirm bool
}

type ConfirmModel struct {
	description string
	value       interface{}
	parent      tea.Model
}

func (c *ConfirmModel) GetDescription() string { return c.description }

func NewConfirmModel(description string, parent tea.Model, value interface{}) (*ConfirmModel, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent is nil")
	}
	return &ConfirmModel{
		description: description,
		parent:      parent,
		value:       value,
	}, nil
}

func (c *ConfirmModel) Init() tea.Cmd { return nil }

func (c *ConfirmModel) View() string {
	var s string

	s += c.description

	base.GetHints(base.EditKey, base.ConfirmKey, base.CancelKey)

	return s
}

func (c *ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case base.ExitKey:
			return c, tea.Quit
		case base.ConfirmKey:
			return c.parent.Update(ConfirmMsg{Value: c.value, Confirm: true})
		case base.CancelKey:
			return c.parent.Update(ConfirmMsg{Value: c.value, Confirm: false})
		}
	}
	return c, nil
}
