package base

type KeyType uint8

const (
	UnknownKeyType KeyType = iota
	UpKeyType
	DownKeyType
	ForwardKeyType
	BackKeyType
	FindKeyType
	EditKeyType
	EnterKeyType
	ExitKeyType
	SelectKeyType
	CancelKeyType
	ConfirmKeyType
	DeleteKeyType
	NextPageKeyType
	PreviousKeyType
)

type ActionKey struct {
	KeyType
	Description string
}
type ActionKeyMap map[string]ActionKey

func (akm ActionKeyMap) SetKey(hotKey string, ak ActionKey) {
	akm[hotKey] = ak
}

func (akm ActionKeyMap) RemoveHotKey(hotKey string) {
	delete(akm, hotKey)
}

// Base KeysMapping
func BaseKeyMappging() ActionKeyMap {
	akm := ActionKeyMap{}
	akm.SetKey("up", ActionKey{KeyType: UpKeyType, Description: "to go up"})
	akm.SetKey("down", ActionKey{KeyType: DownKeyType, Description: "to go down"})
	akm.SetKey("right", ActionKey{KeyType: ForwardKeyType, Description: "to forward"})
	akm.SetKey("left", ActionKey{KeyType: BackKeyType, Description: "to go back"})
	akm.SetKey("f", ActionKey{KeyType: FindKeyType, Description: "to find"})
	akm.SetKey("ctrl+e", ActionKey{KeyType: EditKeyType, Description: "to edit"})
	akm.SetKey("enter", ActionKey{KeyType: EnterKeyType, Description: "to enter"})
	akm.SetKey("ctrl+c", ActionKey{KeyType: ExitKeyType, Description: "to exit"})
	akm.SetKey(" ", ActionKey{KeyType: SelectKeyType, Description: "to select"})
	akm.SetKey("esc", ActionKey{KeyType: CancelKeyType, Description: "to cancel"})
	akm.SetKey("y", ActionKey{KeyType: ConfirmKeyType, Description: "to confirm"})
	akm.SetKey("backspace", ActionKey{KeyType: DeleteKeyType, Description: "to delete"})
	return akm
}

const (
	UpKey      = "up"
	DownKey    = "down"
	ForwardKey = "right"
	BackKey    = "left"
	FindKey    = "f"
	EditKey    = "ctrl+e"
	EnterKey   = "enter"
	ExitKey    = "ctrl+c"
	SelectKey  = " "
	CancelKey  = "esc"
	ConfirmKey = "y"
	DeleteKey  = "backspace"
)
