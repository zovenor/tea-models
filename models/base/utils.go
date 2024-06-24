package base

import (
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

func RepeatSymbol(symbol string, amount int) string {
	var s string
	for i := 0; i < amount; i++ {
		s += symbol
	}
	return s
}

func IsForwardType(t reflect.Value) (tea.Model, bool) {
	if model, ok := t.Interface().(tea.Model); ok {
		return model, true
	} else {
		return nil, false
	}
}
