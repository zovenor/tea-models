package base

import "fmt"

func RepeatSymbol(symbol string, amount int) string {
	var s string
	for i := 0; i < amount; i++ {
		s += symbol
	}
	return s
}

func GetHints(keys ...interface{}) string {
	var s string

	s += "\n"

	for _, key := range keys {
		var k, d string
		switch key := key.(type) {
		case string:
			switch key {
			case UpKey:
				d = "to go up"
			case DownKey:
				d = "to go down"
			case ForwardKey:
				d = "to forward"
			case BackKey:
				d = "to go back"
			case FindKey:
				d = "to find"
			case EditKey:
				d = "to edit"
			case EnterKey:
				d = "to enter"
			case ExitKey:
				d = "to exit the program"
			case SelectKey:
				d = "to select item"
				k = "space"
			case CancelKey:
				d = "to cancel"
			case ConfirmKey:
				d = "to confirm"
			}
			if k == "" {
				k = key
			}
		case Key:
			k = key.Name
			d = key.Description
		}
		s += fmt.Sprintf("Press \"%v\" %v.\n", k, d)
	}
	s += "\n"

	return s
}
