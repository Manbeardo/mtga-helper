package mtga

import (
	"testing"

	"github.com/Manbeardo/mtga-helper/server/mtga/formats"
)

func TestSimulateEvents(t *testing.T) {
	cfg := EventSimulationConfig{
		Format:                  formats.PremierDraft(),
		GameWinRate:             0.567,
		PerEventWinRateVariance: 0.05,
	}
	result := SimulateEvents(cfg, 100000)

	t.Errorf("%#v", result.EconomyStats())
}
