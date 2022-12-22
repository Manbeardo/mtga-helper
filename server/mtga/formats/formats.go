package formats

type MatchKind string

const (
	BestOfOne   = "BO1"
	BestOfThree = "BO3"
)

type Format interface {
	Name() string
	EntryFeeGems() int
	PacksOpened() int
	MatchKind() MatchKind
	MaxLosses() int
	Prizes() []Prize
}

type Prize interface {
	Gems() int
	Packs() int
}

func PremierDraft() Format {
	return format{
		name:         "Premier Draft",
		entryFeeGems: 1500,
		packsOpened:  3,
		matchKind:    BestOfOne,
		maxLosses:    3,
		prizes: []Prize{
			prize{gems: 50, packs: 1},
			prize{gems: 100, packs: 1},
			prize{gems: 250, packs: 2},
			prize{gems: 1000, packs: 2},
			prize{gems: 1400, packs: 3},
			prize{gems: 1600, packs: 4},
			prize{gems: 1800, packs: 5},
			prize{gems: 2200, packs: 6},
		},
	}
}

type format struct {
	name         string
	entryFeeGems int
	packsOpened  int
	matchKind    MatchKind
	maxLosses    int
	prizes       []Prize
}

func (f format) Name() string {
	return f.name
}

func (f format) EntryFeeGems() int {
	return f.entryFeeGems
}

func (f format) PacksOpened() int {
	return f.packsOpened
}

func (f format) MatchKind() MatchKind {
	return f.matchKind
}

func (f format) MaxLosses() int {
	return f.maxLosses
}

func (f format) Prizes() []Prize {
	return f.prizes
}

type prize struct {
	gems  int
	packs int
}

func (p prize) Gems() int {
	return p.gems
}

func (p prize) Packs() int {
	return p.packs
}
