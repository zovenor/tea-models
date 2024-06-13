package base

import "testing"

func TestAddHotKey(t *testing.T) {
	bk := GetBaseKeys()
	err := bk.AddHotKey("f", ExitKeyType)
	if err != nil {
		t.Fatal(err)
	}

}
