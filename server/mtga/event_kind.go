package mtga

import (
	"github.com/Manbeardo/mtga-helper/server/mtga/formats"
	"github.com/Manbeardo/mtga-helper/server/mtga/sets"
)

type EventKind struct {
	sets.Set
	formats.Format
}
