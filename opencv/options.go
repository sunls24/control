package opencv

type findConfig struct {
	scale       int
	sourceRange Range
}

func defFindConfig() *findConfig {
	return &findConfig{
		scale: 1,
	}
}

type FindOption func(*findConfig)

func WithScale(s int) FindOption {
	return func(cfg *findConfig) {
		cfg.scale = s
	}
}

type Range int

const (
	All Range = iota
	Top
	Center
	Bottom
)

func WithRange(sr Range) FindOption {
	return func(cfg *findConfig) {
		cfg.sourceRange = sr
	}
}

func (sr Range) get(rows int) (int, int) {
	switch sr {
	case Top:
		return 0, rows / 3
	case Center:
		return rows / 3, rows * 2 / 3
	case Bottom:
		return rows * 2 / 3, rows
	default:
		return 0, rows
	}
}
