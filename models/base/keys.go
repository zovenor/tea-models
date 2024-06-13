package base

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zovenor/logging/v2"
)

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
	PreviousPageKeyType
)

type ActionKey struct {
	HotKeys     []HotKey
	Description string
	KeyType
}
type ActionKeys []ActionKey

type HotKey string

func NewHotKey(key string, view string) HotKey {
	return HotKey(fmt.Sprintf("%v::%v", key, view))
}
func (hk HotKey) Key() string {
	items := strings.Split(string(hk), "::")
	if len(items) == 0 {
		return ""
	} else {
		return items[0]
	}
}

func (hk HotKey) View() string {
	items := strings.Split(string(hk), "::")
	if len(items) < 2 {
		return hk.Key()
	} else {
		return items[1]
	}
}

func (aks *ActionKeys) AddActionKey(description string, kt KeyType, hotKeys ...HotKey) {
	*aks = append(*aks, ActionKey{HotKeys: hotKeys, Description: description, KeyType: kt})
}

func (aks ActionKeys) AddHotKey(hotKey HotKey, kt KeyType) error {
	for _, ak := range aks {
		for _, hk := range ak.HotKeys {
			if hk.Key() == hotKey.Key() {
				return fmt.Errorf("hotkey already exists in keyType: %v", ak.KeyType)
			}
		}
	}
	for _, ak := range aks {
		if ak.KeyType == kt {
			ak.HotKeys = append(ak.HotKeys, hotKey)
			return nil
		}
	}
	return fmt.Errorf("can not find action key with key type: %v", kt)
}
func (aks ActionKeys) GetActionKeyByKeyType(kt KeyType) (ActionKey, error) {
	for _, ak := range aks {
		if ak.KeyType == kt {
			return ak, nil
		}
	}
	return ActionKey{}, errors.New("can not find action key with key type: " + string(kt))
}
func (aks ActionKeys) GetActionKeyByHotKeyString(hotKey string) (ActionKey, error) {
	for _, ak := range aks {
		for _, hk := range ak.HotKeys {
			if hk.Key() == hotKey {
				return ak, nil
			}
		}
	}
	return ActionKey{}, errors.New("can not find action key with hotkey: " + hotKey)
}
func (aks ActionKeys) GetKeyTypeByHotKeyString(hotKey string) KeyType {
	for _, ak := range aks {
		for _, hk := range ak.HotKeys {
			if hk.Key() == hotKey {
				return ak.KeyType
			}
		}
	}
	return UnknownKeyType
}
func GetBaseKeys() ActionKeys {
	aks := ActionKeys{}
	aks.AddActionKey("to go up", UpKeyType, "up")
	aks.AddActionKey("to go down", DownKeyType, "down")
	aks.AddActionKey("to forward", ForwardKeyType, "right")
	aks.AddActionKey("to go back", BackKeyType, "left")
	aks.AddActionKey("to find", FindKeyType, "f")
	aks.AddActionKey("to edit", EditKeyType, "e")
	aks.AddActionKey("to enter", EnterKeyType, "enter")
	aks.AddActionKey("to exit", ExitKeyType, "ctrl+c")
	aks.AddActionKey("to select", SelectKeyType, " ::space")
	aks.AddActionKey("to cancel", CancelKeyType, "esc")
	aks.AddActionKey("to confirm", ConfirmKeyType, "y")
	aks.AddActionKey("to delete", DeleteKeyType, "backspace")
	aks.AddActionKey("to go to next page", NextPageKeyType, "]")
	aks.AddActionKey("to go to previous page", PreviousPageKeyType, "[")
	return aks
}

func (aks ActionKeys) GetBaseHints(selectedKeyTypes ...KeyType) string {
	var s string

	for _, kt := range selectedKeyTypes {
		if actionKey, err := aks.GetActionKeyByKeyType(kt); err == nil {
			hotKeysString := make([]string, len(actionKey.HotKeys))
			for i, s := range actionKey.HotKeys {
				hotKeysString[i] = string(s.View())
			}
			hotKeysView := strings.Join(hotKeysString, ", ")
			s += fmt.Sprintf("Press %v %v.\n", hotKeysView, actionKey.Description)
		} else {
			logging.FatalSave(err)
		}
	}
	return s
}
