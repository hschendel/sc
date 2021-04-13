package blokus

type Position struct {
	X uint8
	Y uint8
}

func PositionsEqual(a, b []Position) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[Position]struct {
		a int
		b int
	}, len(a))
	for _, p := range a {
		v := m[p]
		v.a++
		m[p] = v
	}
	for _, p := range b {
		v := m[p]
		v.b++
		m[p] = v
	}
	for _, v := range m {
		if v.a != v.b {
			return false
		}
	}
	return true
}
