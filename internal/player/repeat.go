package player

type RepeatMode int

const (
	RepeatOff RepeatMode = iota
	RepeatAll
	RepeatOne
)

func (m RepeatMode) Label() string {
	switch m {
	case RepeatAll:
		return "all"
	case RepeatOne:
		return "one"
	default:
		return "off"
	}
}

func (m RepeatMode) Cycle() RepeatMode {
	return RepeatMode((int(m) + 1) % 3)
}
