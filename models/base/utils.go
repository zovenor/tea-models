package base

func RepeatSymbol(symbol string, amount int) string {
	var s string
	for i := 0; i < amount; i++ {
		s += symbol
	}
	return s
}
